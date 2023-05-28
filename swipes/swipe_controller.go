package swipes

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/github.com/vido21/dating-app/common"
	"github.com/github.com/vido21/dating-app/database"
	premiumPackages "github.com/github.com/vido21/dating-app/premium-packages"
	"github.com/github.com/vido21/dating-app/profiles"
	"github.com/github.com/vido21/dating-app/purchases"
	"github.com/github.com/vido21/dating-app/swipes/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type (
	SwipesController struct {
	}

	SwipeRequest struct {
		Type          string `json:"type" validate:"required"`
		ProfileUserID string `json:"profile_user_id"`
	}
)

func (controller SwipesController) Routes() []common.Route {
	return []common.Route{
		{
			Method:     echo.POST,
			Path:       "/swipe",
			Handler:    controller.Swipe,
			Middleware: []echo.MiddlewareFunc{common.JwtMiddleWare()},
		},
	}
}

func (controller SwipesController) Swipe(ctx echo.Context) error {
	param := new(SwipeRequest)
	if err := ctx.Bind(param); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(param); err != nil && len(param.Type) > 0 {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	var userID uuid.UUID
	// Set default param type to pass
	if len(param.Type) > 0 {
		param.Type = models.Pass
	}
	if len(param.ProfileUserID) > 0 {
		userID = uuid.Must(uuid.FromString(param.ProfileUserID))
	}

	token := ctx.Get("user").(*jwt.Token)
	user := token.Claims.(*common.JwtCustomClaims)

	purchasedPackaged, err := purchases.GetPurchaseService().FindPurchasePackagedByUserID(user.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	// check user is premium or not (from purchase and premium package)
	isUnlimitedQuota := premiumPackages.GetPremiumPackageService().IsConsistsUnlimitedQuotaPackage(purchasedPackaged.PremiumPackages)

	var swipeHistory []models.Swipe
	var excludeUserIDs = []uuid.UUID{
		user.Id,
	}

	if len(userID) > 0 {
		excludeUserIDs = append(excludeUserIDs, userID)
	}

	// get list profile id in swipe history (today)
	db := database.GetInstance()
	err = db.Where("user_id = ?", user.Id).Find(&swipeHistory).Error
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	for _, swipeProfile := range swipeHistory {
		excludeUserIDs = append(excludeUserIDs, swipeProfile.UserID)
	}

	// create swipe history
	if len(userID) > 0 {
		go func(param SwipeRequest, db *gorm.DB) {
			err = db.Create(&models.Swipe{
				SwipeType: param.Type,
				UserID:    user.Id,
				ProfileID: userID,
			}).Error
		}(*param, db)
	}

	if isUnlimitedQuota || len(swipeHistory) < 10 {
		// search profile recomendation first row with condition id not in (user_id, list of swipe history user id)
		profileReccomendation, err := profiles.GetProfileService().GetProfileRecomendation(excludeUserIDs)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, err)

		}
		return ctx.JSON(http.StatusOK, &profileReccomendation)
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "You have reached maximum swipe today",
	})
}
