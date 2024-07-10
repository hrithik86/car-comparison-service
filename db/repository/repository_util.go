package repository

import (
	"car-comparison-service/appcontext/database"
	"gorm.io/gorm"
	"sync"
)

var dbInstance CarComparisonServiceDb
var dbInstanceDoOnce sync.Once

type CarComparisonServiceDb struct {
	*gorm.DB
}

func DbClient() CarComparisonServiceDb {
	dbInstanceDoOnce.Do(func() {
		dbInstance = CarComparisonServiceDb{
			DB: database.GetDB(),
		}
	})

	return dbInstance
}
