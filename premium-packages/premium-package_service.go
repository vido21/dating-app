package premium_packages

import (
	"sync"

	"github.com/github.com/vido21/dating-app/database"
	"github.com/github.com/vido21/dating-app/premium-packages/models"
	uuid "github.com/satori/go.uuid"
)

type premiumPackageService struct{}

var singleton PremiumPackageService
var once sync.Once

func GetPremiumPackageService() PremiumPackageService {
	if singleton != nil {
		return singleton
	}
	once.Do(func() {
		singleton = &premiumPackageService{}
	})
	return singleton
}

func SetPremiumPackageService(service PremiumPackageService) PremiumPackageService {
	original := singleton
	singleton = service
	return original
}

type PremiumPackageService interface {
	FindPremiumPackageByID(id uuid.UUID) *models.PremiumPackage
}

func (u *premiumPackageService) FindPremiumPackageByID(id uuid.UUID) *models.PremiumPackage {
	db := database.GetInstance()

	var premiumPackage models.PremiumPackage
	err := db.First(&premiumPackage, "id = ?", id).Error
	if err == nil {
		return &premiumPackage
	}
	return nil
}
