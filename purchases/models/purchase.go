package models

import (
	"github.com/github.com/vido21/dating-app/common/models"
	PremiumModels "github.com/github.com/vido21/dating-app/premium-packages/models"
	UserModels "github.com/github.com/vido21/dating-app/users/models"
	uuid "github.com/satori/go.uuid"
)

type Purchase struct {
	models.Base
	UserID           uuid.UUID                      `json:"user_id"`
	PremiumPackageID uuid.UUID                      `json:"premium_package_id"`
	PremiumPackages  []PremiumModels.PremiumPackage `json:"premium_packages"`
	Users            []UserModels.User              `json:"users"`
}
