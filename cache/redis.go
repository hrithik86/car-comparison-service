package cache

import (
	"car-comparison-service/config"
	"car-comparison-service/logger"
	"context"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func SetupRedisCluster() error {
	ctx := context.Background()
	redisConf := config.RedisConf()

	client := redis.NewClient(&redis.Options{
		Addr:         redisConf.Host,
		ReadTimeout:  redisConf.ReadTimeout,
		WriteTimeout: redisConf.WriteTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Get(ctx).Error("SetupRedisCluster()", "REDIS CONNECTION ERROR - ", err)
		return err
	}

	redisClient = client

	return nil
}

func GetRedis() *redis.Client {
	return redisClient
}
