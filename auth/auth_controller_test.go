package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/github.com/vido21/dating-app/common"
	"github.com/github.com/vido21/dating-app/common/utils"
	"github.com/github.com/vido21/dating-app/config"
	MocksUtils "github.com/github.com/vido21/dating-app/mocks/common/utils"
	MocksUsers "github.com/github.com/vido21/dating-app/mocks/users"
	commonTest "github.com/github.com/vido21/dating-app/test"
	"github.com/github.com/vido21/dating-app/users"
	UserModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestLoginFailWithInvalidPayload(t *testing.T) {
	println("Login api should return 400 error when the request payload is invalid")
	testServer := echo.New()
	authController := AuthController{}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("invalid json format"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	context := testServer.NewContext(req, resp)

	httpError := authController.Login(context).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
}

func TestLoginFailWithParameterValidation(t *testing.T) {
	println("Login api should return 400 error when the request parameters are invalid")
	testServer := echo.New()
	authController := AuthController{}
	var loginForm LoginRequest
	loginForm.Email = "non-email-format"
	loginForm.Password = "password"
	data, _ := json.Marshal(loginForm)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	context := testServer.NewContext(req, resp)
	httpError := authController.Login(context).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
}

func TestLoginFailWithNonExistingUser(t *testing.T) {
	println("Login api should return 401 error when requested email was not found")
	testServer := echo.New()
	testServer.Validator = &common.CustomValidator{Validator: validator.New()}
	authController := AuthController{}
	loginForm := LoginRequest{
		Email:    commonTest.TestEmail,
		Password: commonTest.TestPassword,
	}
	mockUserService := &MocksUsers.UsersService{}
	mockUserService.On("FindUserByEmail", commonTest.TestEmail).Return(nil)
	originalUserService := users.SetUsersService(mockUserService)

	data, _ := json.Marshal(loginForm)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	context := testServer.NewContext(req, resp)
	httpError := authController.Login(context).(*echo.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, httpError.Code)
	users.SetUsersService(originalUserService)
}

func TestLoginFailWithInvalidPassword(t *testing.T) {
	println("Login api should return 401 error when the password does not match")
	testServer := echo.New()
	testServer.Validator = &common.CustomValidator{Validator: validator.New()}
	authController := AuthController{}
	var loginForm LoginRequest
	loginForm.Email = commonTest.TestEmail
	loginForm.Password = "wrong password"
	mockUserService := &MocksUsers.UsersService{}
	mockUserService.On("FindUserByEmail", commonTest.TestEmail).Return(&UserModels.User{
		Name:     commonTest.TestEmail,
		Password: commonTest.TestPassword,
	})
	originalUserService := users.SetUsersService(mockUserService)

	data, _ := json.Marshal(loginForm)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	context := testServer.NewContext(req, resp)

	httpError := authController.Login(context).(*echo.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, httpError.Code)
	users.SetUsersService(originalUserService)
}

func TestLoginSuccess(t *testing.T) {
	println("Login api should return 200 response when login was successful")
	// create a test user
	user := UserModels.User{
		Name:     commonTest.TestName,
		Email:    commonTest.TestEmail,
		Password: commonTest.TestPassword,
	}
	mockUserService := &MocksUsers.UsersService{}
	mockUserService.On("FindUserByEmail", commonTest.TestEmail).Return(&user)
	originalUserService := users.SetUsersService(mockUserService)

	mockPasswordUtil := MocksUtils.PasswordUtil{}
	mockPasswordUtil.On("CheckPasswordHash", commonTest.TestPassword, commonTest.TestPassword).Return(true)
	originalPasswordUtil := utils.SetPasswordUtil(&mockPasswordUtil)

	testServer := echo.New()
	testServer.Validator = &common.CustomValidator{Validator: validator.New()}
	authController := AuthController{}
	var loginForm LoginRequest
	loginForm.Email = commonTest.TestEmail
	loginForm.Password = commonTest.TestPassword
	data, _ := json.Marshal(loginForm)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	context := testServer.NewContext(req, resp)

	if assert.NoError(t, authController.Login(context)) {
		assert.Equal(t, http.StatusOK, resp.Code)
	}
	users.SetUsersService(originalUserService)
	utils.SetPasswordUtil(originalPasswordUtil)
}

func TestRegisterInvalidPayload(t *testing.T) {
	println("Register api should return 400 error when requested payload is invalid")
	testServer := echo.New()
	testServer.Validator = &common.CustomValidator{Validator: validator.New()}
	authController := AuthController{}
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("invalid json format"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	context := testServer.NewContext(req, resp)

	httpError := authController.Register(context).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
}

func TestRegisterInvalidParams(t *testing.T) {
	println("Register api should return 400 error when requested parameters are invalid")
	testServer := echo.New()
	testServer.Validator = &common.CustomValidator{Validator: validator.New()}
	authController := AuthController{}
	var registerForm RegisterUserRequest
	registerForm.Email = "wrong-email-format"
	registerForm.Password = commonTest.TestPassword
	registerForm.Name = commonTest.TestName
	data, _ := json.Marshal(registerForm)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	context := testServer.NewContext(req, resp)

	httpError := authController.Register(context).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
}

func TestRegisterEmailConflict(t *testing.T) {
	println("Register api should return 400 error when the email is already used")
	testServer := echo.New()
	testServer.Validator = &common.CustomValidator{Validator: validator.New()}
	authController := AuthController{}
	registerForm := RegisterUserRequest{
		Email:    commonTest.TestEmail,
		Password: commonTest.TestPassword,
		Name:     commonTest.TestName,
	}
	mockUserService := &MocksUsers.UsersService{}
	mockUserService.On("FindUserByEmail", commonTest.TestEmail).Return(&UserModels.User{
		Name:     commonTest.TestEmail,
		Password: commonTest.TestPassword,
	})
	originalUserService := users.SetUsersService(mockUserService)
	data, _ := json.Marshal(registerForm)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	context := testServer.NewContext(req, resp)

	httpError := authController.Register(context).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
	users.SetUsersService(originalUserService)
}

func TestRegisterSuccess(t *testing.T) {
	println("Register api should return 200 response when the register is succeeded")
	testServer := echo.New()
	testServer.Validator = &common.CustomValidator{Validator: validator.New()}
	authController := AuthController{}
	registerForm := RegisterUserRequest{
		Email:    commonTest.TestEmail,
		Password: commonTest.TestPassword,
		Name:     commonTest.TestName,
	}
	mockUserService := &MocksUsers.UsersService{}
	mockUserService.On("FindUserByEmail", commonTest.TestEmail).Return(nil)
	mockUserService.On("AddUser", commonTest.TestName, commonTest.TestEmail, commonTest.TestPassword).Return(&UserModels.User{
		Name:     commonTest.TestName,
		Email:    commonTest.TestEmail,
		Password: commonTest.TestPassword,
	})
	originalUserService := users.SetUsersService(mockUserService)
	data, _ := json.Marshal(registerForm)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	context := testServer.NewContext(req, resp)
	if assert.NoError(t, authController.Register(context)) {
		assert.Equal(t, http.StatusOK, resp.Code)
	}
	users.SetUsersService(originalUserService)
}

func TestRoutes(t *testing.T) {
	println("There should be 3 routes defined")
	authController := AuthController{}
	routes := authController.Routes()
	assert.Equal(t, len(routes), 3)
}

func TestProfile(t *testing.T) {
	println("Profile api should return 200 response when the authorization header is valid")
	token, _ := GetAuthService().GetAccessToken(&UserModels.User{
		Name:     commonTest.TestName,
		Email:    commonTest.TestEmail,
		Password: commonTest.TestPassword,
	})

	testServer := echo.New()
	testServer.Use(common.JwtMiddleWare())
	authController := AuthController{}
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
	if assert.NoError(t, authController.Profile(context)) {
		assert.Equal(t, http.StatusOK, resp.Code)
	}
}
