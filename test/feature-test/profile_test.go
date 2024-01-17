package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/github.com/vido21/dating-app/auth"
	"github.com/github.com/vido21/dating-app/common"
	"github.com/github.com/vido21/dating-app/config"
	"github.com/github.com/vido21/dating-app/database"
	profile "github.com/github.com/vido21/dating-app/profiles"
	"github.com/github.com/vido21/dating-app/profiles/models"
	commonTest "github.com/github.com/vido21/dating-app/test"
	userModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateProfileFeature(t *testing.T) {
	commonTest.InitTest()
	t.Run("Register api should return 200 response when the create profile is succeeded", func(t *testing.T) {
		token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
			Name:     commonTest.TestName,
			Email:    commonTest.TestEmail,
			Password: commonTest.TestPassword,
		})

		testServer := echo.New()
		testServer.Validator = &common.CustomValidator{Validator: validator.New()}
		profileController := profile.ProfilesController{}
		form := profile.CreateProfileRequest{
			Sex:            "MALE",
			ProfilePicture: "https://image.com/1",
			About:          "About",
		}

		data, _ := json.Marshal(form)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
		req.Header.Set(echo.HeaderAuthorization, token)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

		if assert.NoError(t, profileController.CreateProfile(context)) {
			assert.Equal(t, http.StatusOK, resp.Code)
		}

	})
}

func TestGetProfileFeature(t *testing.T) {
	commonTest.InitTest()
	t.Run("Profile api should return 200 response when the authorization header is valid", func(t *testing.T) {

		db := database.GetInstance()
		userID, _ := uuid.NewV1()
		var profileModels models.Profile
		profileModels.ProfilePicture = "https://image.com/1"
		profileModels.Sex = 0
		profileModels.About = "About"
		profileModels.UserID = userID
		db.Create(&profileModels)

		token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
			Name:     commonTest.TestName,
			Email:    commonTest.TestEmail,
			Password: commonTest.TestPassword,
		})

		testServer := echo.New()
		testServer.Use(common.JwtMiddleWare())
		authController := profile.ProfilesController{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, token)
		resp := httptest.NewRecorder()
		context := testServer.NewContext(req, resp)
		jwtClaims := common.JwtCustomClaims{
			Name: commonTest.TestName,
			Id:   userID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * config.TokenExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}
		context.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims))

		if assert.NoError(t, authController.GetProfile(context)) {
			assert.Equal(t, http.StatusOK, resp.Code)
		}
	})
}
