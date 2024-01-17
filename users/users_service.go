package users

import (
	"sync"

	"github.com/github.com/vido21/dating-app/database"
	"github.com/github.com/vido21/dating-app/users/models"
)

type usersService struct{}

var singleton UsersService
var once sync.Once

func GetUsersService() UsersService {
	if singleton != nil {
		return singleton
	}
	once.Do(func() {
		singleton = &usersService{}
	})
	return singleton
}

func SetUsersService(service UsersService) UsersService {
	original := singleton
	singleton = service
	return original
}

type UsersService interface {
	FindUserByEmail(email string) *models.User
	AddUser(name string, email string, password string) *models.User
}

func (u *usersService) FindUserByEmail(email string) *models.User {
	db := database.GetInstance()
	var user models.User
	err := db.First(&user, "email = ?", email).Error
	if err == nil {
		return &user
	}
	return nil
}

func (u *usersService) AddUser(name string, email string, password string) *models.User {
	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	db := database.GetInstance()
	db.Create(&user)
	return &user
}
