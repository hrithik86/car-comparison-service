package cache

import (
	"car-comparison-service/config"
	"car-comparison-service/logger"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func SetupRedisCluster() error {
	ctx := context.Background()
	redisConf := config.RedisConf()
	redisUrl := fmt.Sprintf("rediss://%s:%s@%s", redisConf.User, redisConf.Password, redisConf.Host)

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		logger.Get(ctx).Error("SetupRedisCluster()", "error - ", err)
		return err
	}
	client := redis.NewClient(opt)
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
