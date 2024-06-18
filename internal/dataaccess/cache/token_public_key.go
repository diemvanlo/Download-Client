package cache

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"goload/internal/utils"
)

type TokenPublicKey interface {
	Get(ctx context.Context, id uint64) (string, error)
	Set(ctx context.Context, id uint64, data string) error
}

type tokenPublicKey struct {
	client Client
	logger *zap.Logger
}

func NewTokenPublicKey(
	client Client,
	logger *zap.Logger,
) TokenPublicKey {
	return &tokenPublicKey{client: client, logger: logger}
}

func (c tokenPublicKey) getTokenPublicKeyCacheKey(id uint64) string {
	return fmt.Sprintf("token_public_key:%d", id)
}

func (c tokenPublicKey) Get(ctx context.Context, id uint64) (string, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.Uint64("id", id))

	cacheKey := c.getTokenPublicKeyCacheKey(id)
	cacheEntry, err := c.client.Get(ctx, cacheKey)
	if err != nil {
		return "", err
	}

	if cacheEntry == nil {
		return "", nil
	}

	publicKey, ok := cacheEntry.(string)
	if !ok {
		logger.Error("cache entry is not of type bytes")
		return "", nil
	}

	return publicKey, nil
}

func (c tokenPublicKey) Set(ctx context.Context, id uint64, data string) error {
	cacheKey := c.getTokenPublicKeyCacheKey(id)
	return c.client.Set(ctx, cacheKey, data, 0)
}