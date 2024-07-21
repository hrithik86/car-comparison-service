package caching

import (
	"car-comparison-service/cache"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type IRedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	SetWithExpiry(ctx context.Context, key string, val interface{}, expiry time.Duration) error
}

type RedisClient struct {
	RedisClient *redis.Client
}

func GetRedisClient() RedisClient {
	return RedisClient{
		RedisClient: cache.GetRedis(),
	}
}

func (r RedisClient) Get(ctx context.Context, key string) (string, error) {
	var result *redis.StringCmd = r.RedisClient.Get(ctx, key)
	if result.Val() == "" {
		return "", errors.New("key Does Not Exist")
	}
	return result.Val(), nil
}

func (r RedisClient) SetWithExpiry(ctx context.Context, key string, val interface{}, expiry time.Duration) error {
	var result = r.RedisClient.Set(ctx, key, val, expiry)
	if result.Err() != nil {
		return errors.New(result.Err().Error())
	}
	return nil
}
