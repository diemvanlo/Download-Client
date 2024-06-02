package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"goload/internal/configs"
	"goload/internal/dataAccess/cache"
	"goload/internal/dataAccess/database"
	"goload/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	rs512KeyPairCount = 2048
)

var (
	errUnexpectedSigningMethod = status.Errorf(codes.Unauthenticated, "unexpected signing method")
	errCannotGetTokenClaim     = status.Errorf(codes.Unauthenticated, "cannot get token's claims")
	errCannotGetTokenKidClaim  = status.Errorf(codes.Unauthenticated, "cannot get token's kid claims")
	errCannotGetTokenSubClaim  = status.Errorf(codes.Unauthenticated, "cannot get token's sub claims")
	errCannotGetTokenExpClaim  = status.Errorf(codes.Unauthenticated, "cannot get token's exp claims")
	errTokenPublicKeyNotFound  = status.Errorf(codes.Unauthenticated, "token public key not found")
	errInvalidToken            = status.Errorf(codes.Unauthenticated, "invalid token")
	errFailedToSignToken       = status.Errorf(codes.Internal, "failed to sign token")
)

type Token interface {
	GetToken(ctx context.Context, accountID uint64) (string, time.Time, error)
	GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error)
	WithDatabase(database database.Database) Token
}

func generateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	privateKeyPair, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return privateKeyPair, nil
}

type token struct {
	accountDataAccessor        database.AccountDataAccessor
	tokenPublicKeyCache        cache.TokenPublicKey
	tokenPublicKeyDataAccessor database.TokenPublicKeyDataAccessor
	expiresIn                  time.Duration
	privateKey                 *rsa.PrivateKey
	tokenPublicKeyID           uint64
	authConfig                 configs.Auth
	logger                     *zap.Logger
}

func pemEncodePublicKey(key *rsa.PublicKey) ([]byte, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}

	return pem.EncodeToMemory(block), nil
}

func NewToken(accountDataAccessor database.AccountDataAccessor,
	tokenPublicKey cache.TokenPublicKey,
	tokenDataAccessor database.TokenPublicKeyDataAccessor,
	auth configs.Auth,
	logger *zap.Logger) (Token, error) {
	expiresIn, err := auth.Token.GetExpiresInDuration()
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to parse expire_in")
		return nil, err
	}

	rsaKeyPair, err := generateRSAKeyPair(rs512KeyPairCount)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to generate rsa key pair")
		return nil, err
	}

	publicKeyBytes, err := pemEncodePublicKey(&rsaKeyPair.PublicKey)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to encode public key in pem format")
		return nil, err
	}

	tokenPublicKeyID, err := tokenDataAccessor.CreatePublicKey(
		context.Background(),
		database.TokenPublicKey{PublicKey: publicKeyBytes},
	)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create public key entry in database")
		return nil, err
	}

	return &token{
		accountDataAccessor:        accountDataAccessor,
		tokenPublicKeyCache:        tokenPublicKey,
		tokenPublicKeyDataAccessor: tokenDataAccessor,
		expiresIn:                  expiresIn,
		privateKey:                 rsaKeyPair,
		tokenPublicKeyID:           tokenPublicKeyID,
		authConfig:                 auth,
		logger:                     logger,
	}, nil
}

func (t token) GetToken(ctx context.Context, accountID uint64) (string, time.Time, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	expireTime := time.Now().Add(t.expiresIn)
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": accountID,
		"exp": expireTime.Unix(),
		"kid": t.tokenPublicKeyID,
	})

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to sign token")
		return "", time.Time{}, errFailedToSignToken
	}

	return tokenString, expireTime, nil
}

func (t token) GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	parsedToken, err := jwt.Parse(token, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodRSA); !ok {
			logger.Error("unexpected signing method")
			return nil, errUnexpectedSigningMethod
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			logger.Error("cannot get token's claim")
			return nil, errCannotGetTokenClaim
		}

		tokenPublicKeyID, ok := claims["kid"].(float64)
		if !ok {
			logger.Error("Cannot get token's kid claim")
			return nil, errCannotGetTokenKidClaim
		}

		return t.getJWTPublicKey(ctx, uint64(tokenPublicKeyID))
	})

	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to parse token")
	}

	if !parsedToken.Valid {
		logger.Error("invalid token")
		return 0, time.Time{}, errInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("Cannot get token's claim")
		return 0, time.Time{}, errCannotGetTokenClaim
	}

	accountId, ok := claims["sub"].(uint64)
	if !ok {
		logger.Error("cannot get token's sub claim")
		return 0, time.Time{}, errCannotGetTokenSubClaim
	}

	expireTimeUnix, ok := claims["exp"].(float64)
	if !ok {
		logger.Error("Cannot get token's exp claim")
		return 0, time.Time{}, errCannotGetTokenExpClaim
	}

	return uint64(accountId), time.Unix(int64(expireTimeUnix), 0), nil
}

func (t token) getJWTPublicKey(ctx context.Context, id uint64) (*rsa.PublicKey, error) {
	logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64("id", id))

	cachedPublicKeyBytes, err := t.tokenPublicKeyCache.Get(ctx, id)
	if err == nil && cachedPublicKeyBytes != nil {
		return jwt.ParseRSAPublicKeyFromPEM(cachedPublicKeyBytes)
	}
	logger.With(zap.Error(err)).Warn("failed to get cached key bytes, will fall back to database")

	tokenPublicKey, err := t.tokenPublicKeyDataAccessor.GetPublicKey(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errTokenPublicKeyNotFound
		}

		logger.With(zap.Error(err)).Error("Cannot get token's public key from database")
		return nil, err
	}

	err = t.tokenPublicKeyCache.Set(ctx, id, tokenPublicKey.PublicKey)
	if err != nil {
		logger.With(zap.Error(err)).Warn("failed to set public key into cache")
	}

	return jwt.ParseRSAPublicKeyFromPEM(tokenPublicKey.PublicKey)
}

func (t token) WithDatabase(database database.Database) Token {
	t.accountDataAccessor = t.accountDataAccessor.WithDatabase(database)
	return t
}
