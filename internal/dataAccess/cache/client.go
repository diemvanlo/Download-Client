package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"goload/internal/configs"
	"goload/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type Client interface {
	Set(ctx context.Context, key string, data any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	AddToSet(ctx context.Context, key string, data ...any) error
	IsDataSet(ctx context.Context, key string, data any) (bool, error)
}

type client struct {
	redisClient *redis.Client
	logger      *zap.Logger
}

func (c client) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data)).
		With(zap.Duration("ttl", ttl))

	if err := c.redisClient.Set(ctx, key, data, ttl).Err(); err != nil {
		logger.With(zap.Error(err)).Error("Failed to set data into cache")
		return status.Errorf(codes.Internal, "failed to set data into cache: %+v")
	}

	return nil
}

func (c client) Get(ctx context.Context, key string) (any, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("key", key))

	data, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}

		logger.With(zap.Error(err)).Error("Failed to get data from cache")
		return nil, status.Errorf(codes.Internal, "failed to get data from cache: %+v", err)
	}
	return data, nil
}

func (c client) AddToSet(ctx context.Context, key string, data ...any) error {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("key", key)).With(zap.Any("data", data))

	if err := c.redisClient.SAdd(ctx, key, data...).Err(); err != nil {
		logger.With(zap.Error(err)).Error("Failed to set data into set inside cache")
		return status.Errorf(codes.Internal, "failed to set data inside cache: %+v", err)
	}

	return nil
}

func (c client) IsDataSet(ctx context.Context, key string, data any) (bool, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("key", key)).With(zap.Any("data", data))

	result, err := c.redisClient.SIsMember(ctx, key, data).Result()
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to check if data is the member of set inside cache")
		return false, status.Errorf(codes.Internal, "failed to check if data is member of set inside cache: %+v", err)
	}

	return result, nil
}

func NewClient(configs configs.Cache, logger *zap.Logger) Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     configs.Address,
		Username: configs.Username,
		Password: configs.Password,
	})

	return &client{
		redisClient: redisClient,
		logger:      logger,
	}
}
