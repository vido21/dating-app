package utils

import (
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type passwordUtil struct{}

var singleton PasswordUtil
var once sync.Once

func GetPasswordUtil() PasswordUtil {
	if singleton != nil {
		return singleton
	}
	once.Do(func() {
		singleton = &passwordUtil{}
	})
	return singleton
}

func SetPasswordUtil(service PasswordUtil) PasswordUtil {
	original := singleton
	singleton = service
	return original
}

type PasswordUtil interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

func (u *passwordUtil) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *passwordUtil) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
