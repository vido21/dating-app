package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/dgrijalva/jwt-go"
	"github.com/github.com/vido21/dating-app/auth"
	"github.com/github.com/vido21/dating-app/common"
	"github.com/github.com/vido21/dating-app/config"
	"github.com/github.com/vido21/dating-app/database"
	premiumPackage "github.com/github.com/vido21/dating-app/premium-packages"
	commonTest "github.com/github.com/vido21/dating-app/test"
	userModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetListOfPremiumPackages(t *testing.T) {
	commonTest.InitTest()
	t.Run("Register api should return 200 response when the create profile is succeeded", func(t *testing.T) {

		db := database.GetInstance()
		database.SeedPremiumPackages(db)

		token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
			Name:     commonTest.TestName,
			Email:    commonTest.TestEmail,
			Password: commonTest.TestPassword,
		})

		testServer := echo.New()
		testServer.Use(common.JwtMiddleWare())

		premiumPackageController := premiumPackage.PremiumPackageController{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, token)
		resp := httptest.NewRecorder()
		context := testServer.NewContext(req, resp)
		uid, _ := uuid.NewV4()
		jwtClaims := common.JwtCustomClaims{
			Name: commonTest.TestName,
			Id:   uid,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * config.TokenExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}
		context.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims))

		if assert.NoError(t, premiumPackageController.GetListOfPremiumPackages(context)) {
			assert.Equal(t, http.StatusOK, resp.Code)
		}

		monkey.UnpatchAll()

	})
}
