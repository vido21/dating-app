package purchases

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/github.com/vido21/dating-app/common"
	"github.com/github.com/vido21/dating-app/database"
	premiumPackage "github.com/github.com/vido21/dating-app/premium-packages"
	"github.com/github.com/vido21/dating-app/purchases/models"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type (
	PurchaseController struct {
	}

	PurchasePackageRequest struct {
		PremiumPackageID uuid.UUID `json:"premium_package_id" validate:"required"`
		PaymentAmount    float64   `json:"payment_amount" validate:"required"`
	}
)

func (controller PurchaseController) Routes() []common.Route {
	return []common.Route{
		{
			Method:     echo.POST,
			Path:       "/purchase",
			Handler:    controller.PurchasePackage,
			Middleware: []echo.MiddlewareFunc{common.JwtMiddleWare()},
		},
		// {
		// 	Method:     echo.GET,
		// 	Path:       "/profile",
		// 	Handler:    controller.GetProfile,
		// 	Middleware: []echo.MiddlewareFunc{common.JwtMiddleWare()},
		// },
		// {
		// 	Method:  echo.GET,
		// 	Path:    "/profile/next",
		// 	Handler: controller.GetMatch,
		// },
	}
}

func (controller PurchaseController) PurchasePackage(ctx echo.Context) error {
	param := new(PurchasePackageRequest)
	if err := ctx.Bind(param); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(param); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	premium := premiumPackage.GetPremiumPackageService().FindPremiumPackageByID(param.PremiumPackageID)
	if premium == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Premium package not found")
	}

	if premium.Price > param.PaymentAmount {
		return echo.NewHTTPError(http.StatusBadRequest, "Payment amount not enough")
	}

	token := ctx.Get("user").(*jwt.Token)
	user := token.Claims.(*common.JwtCustomClaims)

	db := database.GetInstance()
	var purchase models.Purchase
	purchase.UserID = user.Id
	purchase.PremiumPackageID = premium.ID
	db.Create(&purchase)

	return ctx.JSON(http.StatusOK, purchase)
}
