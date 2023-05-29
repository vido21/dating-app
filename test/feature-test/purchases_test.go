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
	premiumPackagesModels "github.com/github.com/vido21/dating-app/premium-packages/models"
	purchase "github.com/github.com/vido21/dating-app/purchases"
	commonTest "github.com/github.com/vido21/dating-app/test"
	userModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestPurchasePackageFeature(t *testing.T) {
	commonTest.InitTest()
	t.Run("Purchase api should return 200 response when the purchase is succeeded", func(t *testing.T) {

		premium := premiumPackagesModels.PremiumPackage{
			Name:  "Premium Package",
			Price: 200000,
		}
		db := database.GetInstance()
		db.Create(&premium)

		token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
			Name:     commonTest.TestName,
			Email:    commonTest.TestEmail,
			Password: commonTest.TestPassword,
		})

		testServer := echo.New()
		testServer.Validator = &common.CustomValidator{Validator: validator.New()}
		profileController := purchase.PurchaseController{}
		form := purchase.PurchasePackageRequest{
			PremiumPackageID: premium.ID,
			PaymentAmount:    200000,
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

		if assert.NoError(t, profileController.PurchasePackage(context)) {
			assert.Equal(t, http.StatusOK, resp.Code)
		}

	})

	t.Run("Purchase api should return 400 response when the purchase payment amount is less than price", func(t *testing.T) {

		premium := premiumPackagesModels.PremiumPackage{
			Name:  "Premium Package",
			Price: 200000,
		}
		db := database.GetInstance()
		db.Create(&premium)

		token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
			Name:     commonTest.TestName,
			Email:    commonTest.TestEmail,
			Password: commonTest.TestPassword,
		})

		testServer := echo.New()
		testServer.Validator = &common.CustomValidator{Validator: validator.New()}
		profileController := purchase.PurchaseController{}
		form := purchase.PurchasePackageRequest{
			PremiumPackageID: premium.ID,
			PaymentAmount:    100000,
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

		httpError := profileController.PurchasePackage(context).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)

	})

	t.Run("Purchase api should return 404 response when the premium package not found", func(t *testing.T) {

		token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
			Name:     commonTest.TestName,
			Email:    commonTest.TestEmail,
			Password: commonTest.TestPassword,
		})

		testServer := echo.New()
		testServer.Validator = &common.CustomValidator{Validator: validator.New()}
		profileController := purchase.PurchaseController{}
		premiumID, _ := uuid.NewV1()
		form := purchase.PurchasePackageRequest{
			PremiumPackageID: premiumID,
			PaymentAmount:    100000,
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

		httpError := profileController.PurchasePackage(context).(*echo.HTTPError)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

	})
}
