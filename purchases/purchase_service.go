package purchases

import (
	"sync"

	"github.com/github.com/vido21/dating-app/database"
	premiumPackages "github.com/github.com/vido21/dating-app/premium-packages"
	"github.com/github.com/vido21/dating-app/purchases/models"
	uuid "github.com/satori/go.uuid"
)

type purchaseService struct{}

var singleton PurchaseService
var once sync.Once

func GetPurchaseService() PurchaseService {
	if singleton != nil {
		return singleton
	}
	once.Do(func() {
		singleton = &purchaseService{}
	})
	return singleton
}

func SetPurchaseService(service PurchaseService) PurchaseService {
	original := singleton
	singleton = service
	return original
}

type PurchaseService interface {
	FindPurchasePackagedByUserID(userID uuid.UUID) (*models.Purchase, error)
}

func (u *purchaseService) FindPurchasePackagedByUserID(userID uuid.UUID) (*models.Purchase, error) {
	db := database.GetInstance()

	var purchaseData []models.Purchase
	var premiumPackageIDs []uuid.UUID

	err := db.Where("user_id = ?", userID).Find(&purchaseData).Error
	if err != nil {
		return nil, err
	}

	for _, purchase := range purchaseData {
		premiumPackageIDs = append(premiumPackageIDs, purchase.PremiumPackageID)
	}

	purchasedPackaged, err := premiumPackages.GetPremiumPackageService().FindPremiumPackageByIDs(premiumPackageIDs)
	if err != nil {
		return nil, err
	}

	return &models.Purchase{
		PremiumPackages: *purchasedPackaged,
	}, nil
}
