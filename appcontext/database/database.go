package database

import (
	"car-comparison-service/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var gormDb *gorm.DB

func SetupDatabase() error {
	dbConf := config.DbConf()

	connString := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%d",
		dbConf.Name, dbConf.Host, dbConf.User, dbConf.Password, dbConf.Port)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connString,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return err
	}

	gormDb = db
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(int(dbConf.MaxIdleConn))
	sqlDb.SetMaxOpenConns(int(dbConf.MaxOpenConn))
	sqlDb.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifeTime) * time.Minute)
	sqlDb.SetConnMaxIdleTime(time.Duration(dbConf.ConnMaxIdleTime) * time.Minute)

	return nil
}

func GetDB() *gorm.DB {
	return gormDb
}
