package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/github.com/vido21/dating-app/auth"
	"github.com/github.com/vido21/dating-app/common"
	commonTest "github.com/github.com/vido21/dating-app/test"
	"github.com/github.com/vido21/dating-app/users"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterFeature(t *testing.T) {
	commonTest.InitTest()
	t.Run("Register api should return 200 response when the register is succeeded", func(t *testing.T) {
		testServer := echo.New()
		testServer.Validator = &common.CustomValidator{Validator: validator.New()}
		authController := auth.AuthController{}
		form := auth.RegisterUserRequest{
			Email:    commonTest.TestEmail,
			Password: commonTest.TestPassword,
			Name:     commonTest.TestName,
		}

		data, _ := json.Marshal(form)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()
		context := testServer.NewContext(req, resp)
		if assert.NoError(t, authController.Register(context)) {
			assert.Equal(t, http.StatusOK, resp.Code)
		}
	})

	t.Run("Register api should return 400 error when the email is already used", func(t *testing.T) {
		testServer := echo.New()
		testServer.Validator = &common.CustomValidator{Validator: validator.New()}
		authController := auth.AuthController{}
		form := auth.RegisterUserRequest{
			Email:    commonTest.TestName,
			Password: commonTest.TestPassword,
			Name:     commonTest.TestName,
		}

		data, _ := json.Marshal(form)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()
		context := testServer.NewContext(req, resp)

		httpError := authController.Register(context).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestLoginFeature(t *testing.T) {
	commonTest.InitTest()
	t.Run("Login api should return 200 response when login was successful", func(t *testing.T) {
		testServer := echo.New()
		testServer.Validator = &common.CustomValidator{Validator: validator.New()}
		authController := auth.AuthController{}

		users.GetUsersService().AddUser(commonTest.TestName, commonTest.TestEmail, commonTest.TestPassword)

		var loginForm auth.LoginRequest
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
	})

}
