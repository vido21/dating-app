package test

import (
	"log"
	"path/filepath"

	"github.com/github.com/vido21/dating-app/database"
	UserModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/joho/godotenv"
)

func LoadTestEnv() error {
	testEnvPath := filepath.Join("./..", ".test.env")
	err := godotenv.Load(testEnvPath)
	if err != nil {
		log.Fatal("failed to load test env config: ", err)
	}
	return err
}

func InitTest() {
	err := LoadTestEnv()
	if err != nil {
		log.Fatal("failed to load test environment: ", err)
	}

	db := database.GetInstance()
	db.DropTable("migrations")
	db.DropTableIfExists(&UserModels.User{})
	m := database.GetMigrations(db)
	err = m.Migrate()
	if err != nil {
		log.Fatal("failed to run db migration: ", err)
	}
}
