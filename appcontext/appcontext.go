package appcontext

import (
	"car-comparison-service/appcontext/database"
	"car-comparison-service/cache"
	"car-comparison-service/db/repository"
	"github.com/redis/go-redis/v9"
)

type AppContext struct {
	dbClient    repository.CarComparisonServiceDb
	redisClient *redis.Client
}

var appContext *AppContext

func Initiate(module string) error {
	err := database.SetupDatabase()
	if err != nil {
		return err
	}

	err = cache.SetupRedisCluster()
	if err != nil {
		return err
	}

	appContext = &AppContext{
		dbClient:    repository.DbClient(),
		redisClient: cache.GetRedis(),
	}

	return nil
}

func GetDbClient() repository.CarComparisonServiceDb {
	return repository.DbClient()
}
