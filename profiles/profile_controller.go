package profiles

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/github.com/vido21/dating-app/common"
	"github.com/github.com/vido21/dating-app/database"
	"github.com/github.com/vido21/dating-app/profiles/models"
	"github.com/labstack/echo/v4"
)

type (
	ProfilesController struct {
	}

	CreateProfileRequest struct {
		Sex            string `json:"sex" validate:"required"`
		ProfilePicture string `json:"profile_picture" validate:"required"`
		About          string `json:"about"`
	}
)

func (controller ProfilesController) Routes() []common.Route {
	return []common.Route{
		{
			Method:     echo.POST,
			Path:       "/profile",
			Handler:    controller.CreateProfile,
			Middleware: []echo.MiddlewareFunc{common.JwtMiddleWare()},
		},
		{
			Method:     echo.GET,
			Path:       "/profile",
			Handler:    controller.GetProfile,
			Middleware: []echo.MiddlewareFunc{common.JwtMiddleWare()},
		},
	}
}

func (controller ProfilesController) CreateProfile(ctx echo.Context) error {
	params := new(CreateProfileRequest)
	if err := ctx.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(params); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	token := ctx.Get("user").(*jwt.Token)
	user := token.Claims.(*common.JwtCustomClaims)

	db := database.GetInstance()
	var profile models.Profile
	profile.ProfilePicture = params.ProfilePicture
	profile.Sex = models.Sex[params.Sex]
	profile.About = params.About
	profile.UserID = user.Id
	db.Create(&profile)

	return ctx.JSON(http.StatusOK, profile)
}

func (controller ProfilesController) GetProfile(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	user := token.Claims.(*common.JwtCustomClaims)

	userID := user.Id

	db := database.GetInstance()
	var profiles []models.Profile
	db.Find(&profiles)

	var profile models.Profile
	err := db.First(&profile, "user_id = ?", userID).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Profile not found")
	}

	return ctx.JSON(http.StatusOK, profiles)
}
