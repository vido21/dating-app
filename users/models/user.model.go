package models

import (
	CommonModels "github.com/github.com/vido21/dating-app/common/models"
	"github.com/github.com/vido21/dating-app/common/utils"
	ProfileModels "github.com/github.com/vido21/dating-app/profiles/models"
)

type User struct {
	CommonModels.Base
	Email    string                `gorm:"type:varchar(100);unique_index" json:"email"`
	Name     string                `json:"name"`
	Password string                `json:"password"`
	Profile  ProfileModels.Profile `json:"profile"`
}

func (user User) String() string {
	return user.Name
}

func (user *User) BeforeSave() (err error) {
	hashed, err := utils.GetPasswordUtil().HashPassword(user.Password)
	user.Password = hashed
	return
}
