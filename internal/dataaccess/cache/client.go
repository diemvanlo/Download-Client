package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"goload/internal/configs"
	"goload/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
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

type redisClient struct {
	redisClient *redis.Client
	logger      *zap.Logger
}

func (c redisClient) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data)).
		With(zap.Duration("ttl", ttl))

	if err := c.redisClient.Set(ctx, key, data, ttl).Err(); err != nil {
		logger.With(zap.Error(err)).Error("Failed to set data into cache")
		return status.Error(codes.Internal, "failed to set data into cache")
	}

	return nil
}

func (c redisClient) Get(ctx context.Context, key string) (any, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("key", key))

	data, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}

		logger.With(zap.Error(err)).Error("Failed to get data from cache")
		return nil, status.Error(codes.Internal, "failed to get data from cache")
	}
	return data, nil
}

func (c redisClient) AddToSet(ctx context.Context, key string, data ...any) error {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("key", key)).With(zap.Any("data", data))

	if err := c.redisClient.SAdd(ctx, key, data...).Err(); err != nil {
		logger.With(zap.Error(err)).Error("Failed to set data into set inside cache")
		return status.Error(codes.Internal, "failed to set data inside cache")
	}

	return nil
}

func (c redisClient) IsDataSet(ctx context.Context, key string, data any) (bool, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("key", key)).With(zap.Any("data", data))

	result, err := c.redisClient.SIsMember(ctx, key, data).Result()
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to check if data is the member of set inside cache")
		return false, status.Error(codes.Internal, "failed to check if data is member of set inside cache")
	}

	return result, nil
}

func NewRedisClient(cacheConfig configs.Cache, logger *zap.Logger) Client {
	return &redisClient{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     cacheConfig.Address,
			Username: cacheConfig.Username,
			Password: cacheConfig.Password,
		}),
		logger: logger,
	}
}

type inMemoryClient struct {
	cache      map[string]any
	cacheMutex *sync.Mutex
	logger     *zap.Logger
}

func (i inMemoryClient) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	i.cache[key] = data
	return nil
}

func (i inMemoryClient) Get(ctx context.Context, key string) (any, error) {
	data, ok := i.cache[key]
	if !ok {
		return nil, ErrCacheMiss
	}

	return data, nil
}

func (i inMemoryClient) AddToSet(ctx context.Context, key string, data ...any) error {
	i.cacheMutex.Lock()
	defer i.cacheMutex.Unlock()
	set := i.getSet(key)
	set = append(set, data...)
	i.cache[key] = set
	return nil
}

func (i inMemoryClient) IsDataSet(ctx context.Context, key string, data any) (bool, error) {
	i.cacheMutex.Lock()
	defer i.cacheMutex.Unlock()

	set := i.getSet(key)

	for i := range set {
		if set[i] == data {
			return true, nil
		}
	}

	return false, nil
}

func (i inMemoryClient) getSet(key string) []any {
	setValue, ok := i.cache[key]
	if !ok {
		return make([]any, 0)
	}

	set, ok := setValue.([]any)
	if !ok {
		return make([]any, 0)
	}

	return set
}

func NewInMemoryClient(logger *zap.Logger) Client {
	return &inMemoryClient{
		cache:      make(map[string]any),
		cacheMutex: new(sync.Mutex),
		logger:     logger,
	}
}

func NewClient(cacheConfigs configs.Cache, logger *zap.Logger) (Client, error) {
	switch cacheConfigs.Type {
	case configs.CacheTypeInMemory:
		return NewInMemoryClient(logger), nil
	case configs.CacheTypeRedis:
		return NewRedisClient(cacheConfigs, logger), nil
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheConfigs.Type)
	}
}
