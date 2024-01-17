package premium_packages

import (
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
	commonTest "github.com/github.com/vido21/dating-app/test"
	userModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	println("There should be 1 route defined")
	authController := PremiumPackageController{}
	routes := authController.Routes()
	assert.Equal(t, len(routes), 1)
}

func TestGetListOfPremiumPackagesController(t *testing.T) {
	tests := []struct {
		name string
		test func()
	}{
		{
			name: "Get List of Premium Package api should return 200 response when the authorization header is valid and success get data",
			test: func() {
				token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
					Name:     commonTest.TestName,
					Email:    commonTest.TestEmail,
					Password: commonTest.TestPassword,
				})

				testServer := echo.New()
				testServer.Use(common.JwtMiddleWare())

				premiumPackageController := PremiumPackageController{}

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
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Find", func(db *gorm.DB, out interface{}, where ...interface{}) *gorm.DB {
					return &gorm.DB{
						Error: nil,
					}
				})

				if assert.NoError(t, premiumPackageController.GetListOfPremiumPackages(context)) {
					assert.Equal(t, http.StatusOK, resp.Code)
				}

				monkey.UnpatchAll()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test()
		})
	}
}
