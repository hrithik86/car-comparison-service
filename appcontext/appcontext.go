package appcontext

import (
	"car-comparison-service/appcontext/database"
	"car-comparison-service/db/repository"
)

type AppContext struct {
	dbClient repository.CarComparisonServiceDb
}

var appContext *AppContext

func Initiate(module string) error {
	err := database.SetupDatabase()
	if err != nil {
		return err
	}

	appContext = &AppContext{
		dbClient: repository.DbClient(),
	}

	return nil
}

func GetDbClient() repository.CarComparisonServiceDb {
	return repository.DbClient()
}
