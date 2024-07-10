package appcontext

import (
	"car-comparison-service/appcontext/database"
)

type AppContext struct {
	dbClient database.CarComparisonServiceDb
}

var appContext *AppContext

func Initiate(module string) error {
	err := database.SetupDatabase()
	if err != nil {
		return err
	}

	appContext = &AppContext{
		dbClient: database.GetDbClient(),
	}

	return nil
}
