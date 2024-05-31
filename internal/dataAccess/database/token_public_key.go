package database

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"goload/internal/utils"
)

const (
	TabNameTokenPublicKeys         = "token_public_keys"
	ColNameTokenPublicKeyID        = "id"
	ColNameTokenPublicKeyPublicKey = "public_key"
)

type TokenPublicKey struct {
	ID        uint64 `sql:"id"`
	PublicKey []byte `sql:"public_key"`
}

type TokenPublicKeyDataAccessor interface {
	CreatePublicKey(ctx context.Context, tokenPublicKey TokenPublicKey) (uint64, error)
	GetPublicKey(ctx context.Context, id uint64) (TokenPublicKey, error)
	WithDatabase(database Database) TokenPublicKeyDataAccessor
}

type tokenPublicKeyDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func (t tokenPublicKeyDataAccessor) CreatePublicKey(ctx context.Context, tokenPublicKey TokenPublicKey) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)
	result, err := t.database.Insert(TabNameTokenPublicKeys).Rows(goqu.Record{
		ColNameTokenPublicKeyPublicKey: tokenPublicKey.PublicKey,
	}).Executor().ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to create token public key")
		return 0, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get last inserted id")
		return 0, err
	}

	return uint64(lastInsertedID), nil
}

func (t tokenPublicKeyDataAccessor) GetPublicKey(ctx context.Context, id uint64) (TokenPublicKey, error) {
	logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64(ColNameTokenPublicKeyID, id))

	tokenPublicKey := TokenPublicKey{}
	found, err := t.database.Select().From(TabNameTokenPublicKeys).Where(goqu.Ex{
		ColNameTokenPublicKeyID: id,
	}).Executor().ScanStructContext(ctx, &tokenPublicKey)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get public key")
		return TokenPublicKey{}, err
	}

	if !found {
		logger.Warn("public key not found")
		return TokenPublicKey{}, sql.ErrNoRows
	}

	return tokenPublicKey, nil
}

func (t tokenPublicKeyDataAccessor) WithDatabase(database Database) TokenPublicKeyDataAccessor {
	t.database = database
	return t
}

func NewTokenPublicKeyDataAccessor(database *goqu.Database, logger *zap.Logger) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		database: database,
		logger:   logger,
	}
}
