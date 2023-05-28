package database

import (
	PremiumPackageModels "github.com/github.com/vido21/dating-app/premium-packages/models"
	ProfileModels "github.com/github.com/vido21/dating-app/profiles/models"
	PurchaseModels "github.com/github.com/vido21/dating-app/purchases/models"
	SwipeModels "github.com/github.com/vido21/dating-app/swipes/models"
	UserModels "github.com/github.com/vido21/dating-app/users/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func GetMigrations(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "2020080202",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&UserModels.User{}).Error; err != nil {
					return err
				}
				if err := tx.AutoMigrate(&ProfileModels.Profile{}).Error; err != nil {
					return err
				}
				if err := tx.AutoMigrate(&PremiumPackageModels.PremiumPackage{}).Error; err != nil {
					return err
				}
				if err := tx.AutoMigrate(&PurchaseModels.Purchase{}).Error; err != nil {
					return err
				}
				if err := tx.AutoMigrate(&SwipeModels.Swipe{}).Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.DropTable("blogs").Error; err != nil {
					return nil
				}
				if err := tx.DropTable("users").Error; err != nil {
					return nil
				}
				if err := tx.DropTable("premium_packages").Error; err != nil {
					return nil
				}
				if err := tx.DropTable("profiles").Error; err != nil {
					return nil
				}
				if err := tx.DropTable("purchases").Error; err != nil {
					return nil
				}
				if err := tx.DropTable("swipes").Error; err != nil {
					return nil
				}
				return nil
			},
		},
	})
}
