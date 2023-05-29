package purchases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	println("There should be 1 route defined")
	authController := PurchaseController{}
	routes := authController.Routes()
	assert.Equal(t, len(routes), 1)
}

// func TestPurchaseController_PurchasePackage(t *testing.T) {
// 	println("Register api should return 200 response when the register is succeeded")
// 	token, _ := auth.GetAuthService().GetAccessToken(&userModels.User{
// 		Name:     commonTest.TestName,
// 		Email:    commonTest.TestEmail,
// 		Password: commonTest.TestPassword,
// 	})

// 	testServer := echo.New()
// 	testServer.Validator = &common.CustomValidator{Validator: validator.New()}
// 	purchaseController := PurchaseController{}

// 	packageID, _ := uuid.NewV1()
// 	registerForm := PurchasePackageRequest{
// 		PremiumPackageID: packageID,
// 		PaymentAmount:    200000,
// 	}

// 	data, _ := json.Marshal(registerForm)
// 	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(data)))
// 	req.Header.Set(echo.HeaderAuthorization, token)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	resp := httptest.NewRecorder()
// 	context := testServer.NewContext(req, resp)

// 	mockUserService := &mockPremiumPackage.PremiumPackageService{}
// 	mockUserService.On("FindPremiumPackageByID", packageID).Return(&premiumPackageModels.PremiumPackage{
// 		Type:  premiumPackageModels.UnilimitedQuota,
// 		Price: 200000,
// 	})
// 	originalPremiumPackageService := premiumPackageService.SetPremiumPackageService(mockUserService)

// 	userID, _ := uuid.NewV1()
// 	monkey.Patch(context.Get, func(key string) interface{} {
// 		return &jwt.Token{
// 			Claims: common.JwtCustomClaims{
// 				Id: userID,
// 			},
// 		}
// 	})

// 	monkey.Patch(database.GetInstance, func() *gorm.DB {
// 		return &gorm.DB{}
// 	})
// 	monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Create", func(db *gorm.DB, value interface{}) *gorm.DB {
// 		return &gorm.DB{
// 			Error: nil,
// 		}
// 	})

// 	if assert.NoError(t, purchaseController.PurchasePackage(context)) {
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 	}
// 	premiumPackageService.SetPremiumPackageService(originalPremiumPackageService)
// 	monkey.UnpatchAll()
// }
