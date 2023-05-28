package test

import (
	"log"
	"os"

	BlogModels "github.com/github.com/vido21/dating-app/blogs/models"
	"github.com/github.com/vido21/dating-app/database"
	UserModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/joho/godotenv"
)

func LoadTestEnv() error {
	err := godotenv.Load(os.ExpandEnv("$GOPATH/src/github.com/github.com/vido21/dating-app/test.env"))
	if err != nil {
		log.Fatal("failed to load test env config: ", err)
	}
	return err
}

func InitTest() {
	err := LoadTestEnv()
	db := database.GetInstance()
	db.DropTable("migrations")
	db.DropTableIfExists(&UserModels.User{})
	db.DropTableIfExists(&BlogModels.Blog{})
	m := database.GetMigrations(db)
	err = m.Migrate()
	if err != nil {
		log.Fatal("failed to run db migration: ", err)
	}
}
