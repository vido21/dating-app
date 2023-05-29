package premium_packages

import (
	"net/http"

	"github.com/github.com/vido21/dating-app/common"
	"github.com/github.com/vido21/dating-app/database"
	"github.com/github.com/vido21/dating-app/premium-packages/models"
	"github.com/labstack/echo/v4"
)

type (
	PremiumPackageController struct {
	}

	CreateProfileRequest struct {
		Sex            string `json:"sex" validate:"required"`
		ProfilePicture string `json:"profile_picture" validate:"required"`
		About          string `json:"about"`
	}
)

func (controller PremiumPackageController) Routes() []common.Route {
	return []common.Route{
		{
			Method:     echo.GET,
			Path:       "/premium-package",
			Handler:    controller.GetListOfPremiumPackages,
			Middleware: []echo.MiddlewareFunc{common.JwtMiddleWare()},
		},
	}
}

func (controller PremiumPackageController) GetListOfPremiumPackages(ctx echo.Context) error {
	db := database.GetInstance()
	var premiumPackages []models.PremiumPackage
	db.Find(&premiumPackages)

	if err := db.Find(&premiumPackages).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Premium Package not found")
	}

	return ctx.JSON(http.StatusOK, premiumPackages)
}
