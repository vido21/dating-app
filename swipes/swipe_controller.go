package swipes

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/github.com/vido21/dating-app/common"
	appError "github.com/github.com/vido21/dating-app/common/error"
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
		log.Println(appError.AddTrace(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(param); err != nil && len(param.Type) > 0 {
		log.Println(appError.AddTrace(err))
		return ctx.JSON(http.StatusBadRequest, err)
	}

	var profileUserID uuid.UUID
	// Set default param type to pass
	if len(param.Type) > 0 {
		param.Type = models.Pass
	}

	if len(param.ProfileUserID) > 0 {
		profileUserID = uuid.Must(uuid.FromString(param.ProfileUserID))
	}

	token := ctx.Get("user").(*jwt.Token)
	user := token.Claims.(*common.JwtCustomClaims)

	purchasedPackaged, err := purchases.GetPurchaseService().FindPurchasePackagedByUserID(user.Id)
	if err != nil {
		log.Println(appError.AddTrace(err))
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	// check user is premium or not (from purchase and premium package)
	isUnlimitedQuota := premiumPackages.GetPremiumPackageService().IsConsistsUnlimitedQuotaPackage(purchasedPackaged.PremiumPackages)

	var swipeHistory []models.Swipe
	var excludeProfileIDs = []uuid.UUID{
		user.Id,
	}

	if profileUserID != uuid.Nil {
		excludeProfileIDs = append(excludeProfileIDs, profileUserID)
	}

	// get list profile id in swipe history (today)
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	db := database.GetInstance()
	err = db.Where("user_id = ? and created_at BETWEEN ? AND ?", user.Id, startOfDay, endOfDay).Find(&swipeHistory).Error
	if err != nil {
		log.Println(appError.AddTrace(err))
		return ctx.JSON(http.StatusNotFound, err)
	}

	for _, swipeProfile := range swipeHistory {
		excludeProfileIDs = append(excludeProfileIDs, swipeProfile.ProfileID)
	}

	// create swipe history
	if len(profileUserID) > 0 {
		go func(param SwipeRequest, db *gorm.DB) {
			err = db.Create(&models.Swipe{
				SwipeType: strings.ToUpper(param.Type),
				UserID:    user.Id,
				ProfileID: profileUserID,
			}).Error
		}(*param, db)

		if err != nil {
			log.Println(appError.AddTrace(err))
		}
	}

	if isUnlimitedQuota || len(swipeHistory) < models.LimitSwipe {
		// search profile recomendation first row with condition id not in (user_id, list of swipe history user id)
		profileReccomendation, err := profiles.GetProfileService().GetProfileRecomendation(excludeProfileIDs, user.Id)
		if err != nil {
			log.Println(appError.AddTrace(err))
			return ctx.JSON(http.StatusBadRequest, err)
		}
		return ctx.JSON(http.StatusOK, &profileReccomendation)
	}

	return ctx.JSON(http.StatusNotAcceptable, map[string]interface{}{
		"message": "You have reached maximum swipe today",
	})
}
