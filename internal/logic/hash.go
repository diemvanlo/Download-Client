package logic

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"goload/internal/configs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Hash interface {
	Hash(ctx context.Context, data string) (string, error)
	IsHashEqual(ctx context.Context, data string, hashed string) (bool, error)
}

type hash struct {
	accountConfig configs.Auth
}

func NewHash(config configs.Auth) Hash {
	return &hash{accountConfig: config}
}

func (h hash) Hash(ctx context.Context, data string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(data), h.accountConfig.Hash.HashCost)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to hash data")
	}

	return string(hashed), nil
}

func (h hash) IsHashEqual(ctx context.Context, data string, hashed string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(data)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, status.Error(codes.Internal, "failed to if data equal hash")
	}

	return true, nil
}
