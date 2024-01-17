package purchases

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/github.com/vido21/dating-app/database"
	premiumPackages "github.com/github.com/vido21/dating-app/premium-packages"
	mockPremiumPackage "github.com/github.com/vido21/dating-app/premium-packages/mocks"
	premiumPackageModels "github.com/github.com/vido21/dating-app/premium-packages/models"
	"github.com/github.com/vido21/dating-app/purchases/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func Test_purchaseService_FindPurchasePackagedByUserID(t *testing.T) {
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		u       *purchaseService
		args    args
		want    *models.Purchase
		wantErr bool
		mock    func()
	}{
		{
			name: "success get data",
			u:    &purchaseService{},
			args: args{
				userID: uuid.UUID{},
			},
			want: &models.Purchase{},
			mock: func() {
				monkey.Patch(database.GetInstance, func() *gorm.DB {
					return &gorm.DB{}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Where", func(db *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
					return &gorm.DB{
						Error: nil,
					}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Find", func(db *gorm.DB, out interface{}, where ...interface{}) *gorm.DB {
					return &gorm.DB{
						Error: nil,
					}
				})

				mockPremium := &mockPremiumPackage.PremiumPackageService{}
				mockPremium.On("FindPremiumPackageByIDs", []uuid.UUID{}).Return(&[]premiumPackageModels.PremiumPackage{}, nil)
				originalPremiumService := premiumPackages.SetPremiumPackageService(mockPremium)
				premiumPackages.SetPremiumPackageService(originalPremiumService)
			},
		},
		{
			name: "error get purchase data",
			u:    &purchaseService{},
			args: args{
				userID: uuid.UUID{},
			},
			want:    nil,
			wantErr: true,
			mock: func() {
				monkey.Patch(database.GetInstance, func() *gorm.DB {
					return &gorm.DB{}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Where", func(db *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
					return &gorm.DB{
						Error: nil,
					}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Find", func(db *gorm.DB, out interface{}, where ...interface{}) *gorm.DB {
					return &gorm.DB{
						Error: gorm.ErrRecordNotFound,
					}
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &purchaseService{}

			tt.mock()
			defer monkey.UnpatchAll()

			got, err := u.FindPurchasePackagedByUserID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("purchaseService.FindPurchasePackagedByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("purchaseService.FindPurchasePackagedByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPurchaseService(t *testing.T) {
	tests := []struct {
		name string
		want PurchaseService
		mock func()
	}{
		{
			name: "singleton is nil",
			want: &purchaseService{},
		},
		{
			name: "singleton is nil",
			want: &purchaseService{},
			mock: func() {
				singleton = &purchaseService{}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.mock != nil {
				tt.mock()
			}

			if got := GetPurchaseService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPurchaseService() = %v, want %v", got, tt.want)
			}
		})
	}
}
