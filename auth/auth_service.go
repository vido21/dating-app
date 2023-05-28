package auth

import (
	"os"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/github.com/vido21/dating-app/common"
	"github.com/github.com/vido21/dating-app/config"
	UserModels "github.com/github.com/vido21/dating-app/users/models"
)

type authService struct{}

var singleton AuthService
var once sync.Once

func GetAuthService() AuthService {
	once.Do(func() {
		singleton = &authService{}
	})
	return singleton
}

//func SetAuthService(service AuthService) AuthService {
//	original := singleton
//	singleton = service
//	return original
//}

type AuthService interface {
	GetAccessToken(user *UserModels.User) (string, error)
}

func (s *authService) GetAccessToken(user *UserModels.User) (string, error) {
	claims := &common.JwtCustomClaims{
		Name:  user.Name,
		Id:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * config.TokenExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
