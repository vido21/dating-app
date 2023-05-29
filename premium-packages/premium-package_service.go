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
	FindPremiumPackageByID(id uuid.UUID) (*models.PremiumPackage, error)
	FindPremiumPackageByIDs(ids []uuid.UUID) (*[]models.PremiumPackage, error)
	IsConsistsUnlimitedQuotaPackage(premiumPackages []models.PremiumPackage) bool
	IsConsistsVerifiedUserPackage(premiumPackages []models.PremiumPackage) bool
}

func (u *premiumPackageService) FindPremiumPackageByID(id uuid.UUID) (*models.PremiumPackage, error) {
	db := database.GetInstance()
	var premiumPackage models.PremiumPackage

	err := db.First(&premiumPackage, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &premiumPackage, nil
}

func (u *premiumPackageService) FindPremiumPackageByIDs(ids []uuid.UUID) (*[]models.PremiumPackage, error) {
	db := database.GetInstance()
	var premiumPackage []models.PremiumPackage

	err := db.Where("id IN (?)", ids).Find(&premiumPackage).Error
	if err != nil {
		return nil, err
	}

	return &premiumPackage, nil
}

func (u *premiumPackageService) IsConsistsUnlimitedQuotaPackage(premiumPackages []models.PremiumPackage) bool {
	for _, premium := range premiumPackages {
		if premium.Type == models.UnilimitedQuota {
			return true
		}
	}

	return false
}

func (u *premiumPackageService) IsConsistsVerifiedUserPackage(premiumPackages []models.PremiumPackage) bool {
	for _, premium := range premiumPackages {
		if premium.Type == models.VerifiedUser {
			return true
		}
	}

	return false
}
