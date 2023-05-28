package database

import (
	premiumPackageModels "github.com/github.com/vido21/dating-app/premium-packages/models"
	"github.com/jinzhu/gorm"
)

func SeedPremiumPackages(db *gorm.DB) error {
	premiumPackages := []premiumPackageModels.PremiumPackage{
		{
			Description: "Upgrade to our Premium Package and enjoy unlimited access to all the features of our dating app. With the Unlimited Quota package, you can swipe right and left without any restrictions on the number of profiles you can view in a day. Find your perfect match without limitations.",
			Name:        "Premium Package - Unlimited Quota",
			Type:        premiumPackageModels.UnilimitedQuota,
			Price:       float64(200000),
		},
		{
			Description: "Upgrade to our Premium Package and get the Verified Label. Stand out from the crowd with a verified badge on your profile, showing others that you're a trusted user. Increase your chances of making meaningful connections and enjoy a more enhanced online dating experience.",
			Name:        "Premium Package - Verified Label",
			Type:        premiumPackageModels.VerifiedUser,
			Price:       float64(100000),
		},
	}

	for _, premiumPackage := range premiumPackages {
		err := db.Where(premiumPackageModels.PremiumPackage{Type: premiumPackage.Type}).FirstOrCreate(&premiumPackage).Error
		if err != nil {
			return err
		}
	}

	return nil
}
