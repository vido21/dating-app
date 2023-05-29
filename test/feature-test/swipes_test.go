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
	profileModels "github.com/github.com/vido21/dating-app/profiles/models"
	swipe "github.com/github.com/vido21/dating-app/swipes"
	"github.com/github.com/vido21/dating-app/swipes/models"
	commonTest "github.com/github.com/vido21/dating-app/test"
	userModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestSwipeFeature(t *testing.T) {
	commonTest.InitTest()
	t.Run("Swipe api should return 200 response and success get data", func(t *testing.T) {

		userID1, _ := uuid.NewV1()
		db := database.GetInstance()
		db.Create(&profileModels.Profile{
			ProfilePicture: "https://image.com/1",
			Sex:            0,
			About:          "About",
			UserID:         userID1,
		})

		db.Create(&profileModels.Profile{
			ProfilePicture: "https://image.com/1",
			Sex:            1,
			About:          "About",
		})

		token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
			Name:     commonTest.TestName,
			Email:    commonTest.TestEmail,
			Password: commonTest.TestPassword,
		})

		testServer := echo.New()
		testServer.Validator = &common.CustomValidator{Validator: validator.New()}
		profileController := swipe.SwipesController{}
		form := swipe.SwipeRequest{
			Type:          "pass",
			ProfileUserID: userID1.String(),
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

		if assert.NoError(t, profileController.Swipe(context)) {
			assert.Equal(t, http.StatusOK, resp.Code)
		}

	})

	t.Run("Swipe api should return 200 response and have reached maximum swipe", func(t *testing.T) {

		userID1, _ := uuid.NewV1()
		db := database.GetInstance()
		db.Create(&profileModels.Profile{
			ProfilePicture: "https://image.com/1",
			Sex:            1,
			About:          "About",
			UserID:         userID1,
		})

		db.Create(&profileModels.Profile{
			ProfilePicture: "https://image.com/1",
			Sex:            1,
			About:          "About",
		})

		for i := 0; i < 10; i++ {
			db.Create(&models.Swipe{
				UserID: userID1,
			})
		}

		token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
			Name:     commonTest.TestName,
			Email:    commonTest.TestEmail,
			Password: commonTest.TestPassword,
		})

		testServer := echo.New()
		testServer.Validator = &common.CustomValidator{Validator: validator.New()}
		profileController := swipe.SwipesController{}
		form := swipe.SwipeRequest{
			Type: "pass",
		}

		data, _ := json.Marshal(form)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
		req.Header.Set(echo.HeaderAuthorization, token)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()
		context := testServer.NewContext(req, resp)
		jwtClaims := common.JwtCustomClaims{
			Name: commonTest.TestName,
			Id:   userID1,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * config.TokenExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}
		context.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims))

		if assert.NoError(t, profileController.Swipe(context)) {
			assert.Equal(t, http.StatusNotAcceptable, resp.Code)
		}

	})
}
