package profiles

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/dgrijalva/jwt-go"
	"github.com/github.com/vido21/dating-app/auth"
	"github.com/github.com/vido21/dating-app/common"
	"github.com/github.com/vido21/dating-app/config"
	"github.com/github.com/vido21/dating-app/database"
	"github.com/github.com/vido21/dating-app/test"
	commonTest "github.com/github.com/vido21/dating-app/test"
	userModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	println("There should be 1 route defined")
	authController := ProfilesController{}
	routes := authController.Routes()
	assert.Equal(t, len(routes), 2)
}

func TestProfilesController_CreateProfile(t *testing.T) {
	println("Register api should return 200 response when the register is succeeded")
	token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
		Name:     commonTest.TestName,
		Email:    commonTest.TestEmail,
		Password: commonTest.TestPassword,
	})

	testServer := echo.New()
	testServer.Validator = &common.CustomValidator{Validator: validator.New()}
	profileController := ProfilesController{}
	registerForm := CreateProfileRequest{
		Sex:            "MALE",
		ProfilePicture: "https://image.com/1",
		About:          "About",
	}

	data, _ := json.Marshal(registerForm)
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

	monkey.Patch(database.GetInstance, func() *gorm.DB {
		return &gorm.DB{}
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Create", func(db *gorm.DB, value interface{}) *gorm.DB {
		return &gorm.DB{
			Error: nil,
		}
	})

	if assert.NoError(t, profileController.CreateProfile(context)) {
		assert.Equal(t, http.StatusOK, resp.Code)
	}

	monkey.UnpatchAll()
}

func TestProfilesController_GetProfile(t *testing.T) {
	println("Profile api should return 200 response when the authorization header is valid")
	test.LoadTestEnv()
	token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
		Name:     commonTest.TestName,
		Email:    commonTest.TestEmail,
		Password: commonTest.TestPassword,
	})

	testServer := echo.New()
	testServer.Use(common.JwtMiddleWare())
	authController := ProfilesController{}
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

	monkey.Patch(database.GetInstance, func() *gorm.DB {
		return &gorm.DB{}
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "First", func(db *gorm.DB, out interface{}, where ...interface{}) *gorm.DB {
		return &gorm.DB{
			Error: nil,
		}
	})

	if assert.NoError(t, authController.GetProfile(context)) {
		assert.Equal(t, http.StatusOK, resp.Code)
	}
	monkey.UnpatchAll()
}
