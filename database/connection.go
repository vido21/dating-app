package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/github.com/vido21/dating-app/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var onceDb sync.Once

var instance *gorm.DB

func GetInstance() *gorm.DB {
	onceDb.Do(func() {
		databaseConfig := config.DatabaseNew().(*config.DatabaseConfig)
		db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			databaseConfig.Psql.DbHost,
			databaseConfig.Psql.DbPort,
			databaseConfig.Psql.DbUsername,
			databaseConfig.Psql.DbDatabase,
			databaseConfig.Psql.DbPassword,
		))
		if err != nil {
			log.Fatalf("Could not connect to database :%v", err)
		}
		instance = db
	})
	return instance
}
