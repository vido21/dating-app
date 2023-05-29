package premium_packages

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/github.com/vido21/dating-app/database"
	"github.com/github.com/vido21/dating-app/premium-packages/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func Test_premiumPackageService_FindPremiumPackageByID(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		u       *premiumPackageService
		args    args
		want    *models.PremiumPackage
		wantErr bool
		mock    func()
	}{
		{
			name: "success get data",
			u:    &premiumPackageService{},
			args: args{
				id: uuid.UUID{},
			},
			want: &models.PremiumPackage{},
			mock: func() {
				monkey.Patch(database.GetInstance, func() *gorm.DB {
					return &gorm.DB{}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "First", func(db *gorm.DB, out interface{}, where ...interface{}) *gorm.DB {
					return &gorm.DB{
						Error: nil,
					}
				})
			},
		},
		{
			name: "error get data",
			u:    &premiumPackageService{},
			args: args{
				id: uuid.UUID{},
			},
			wantErr: true,
			mock: func() {
				monkey.Patch(database.GetInstance, func() *gorm.DB {
					return &gorm.DB{}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "First", func(db *gorm.DB, out interface{}, where ...interface{}) *gorm.DB {
					return &gorm.DB{
						Error: gorm.ErrRecordNotFound,
					}
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &premiumPackageService{}

			tt.mock()
			defer monkey.UnpatchAll()

			got, err := u.FindPremiumPackageByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("premiumPackageService.FindPremiumPackageByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_premiumPackageService_FindPremiumPackageByIDs(t *testing.T) {
	type args struct {
		ids []uuid.UUID
	}

	var premiumPackage []models.PremiumPackage

	tests := []struct {
		name    string
		u       *premiumPackageService
		args    args
		want    *[]models.PremiumPackage
		wantErr bool
		mock    func()
	}{
		{
			name: "success get data",
			u:    &premiumPackageService{},
			args: args{
				ids: []uuid.UUID{},
			},
			want: &premiumPackage,
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
			},
		},
		{
			name: "error get data",
			u:    &premiumPackageService{},
			args: args{
				ids: []uuid.UUID{},
			},
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
			u := &premiumPackageService{}

			tt.mock()
			defer monkey.UnpatchAll()

			got, err := u.FindPremiumPackageByIDs(tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("premiumPackageService.FindPremiumPackageByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("premiumPackageService.FindPremiumPackageByIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_premiumPackageService_IsConsistsUnlimitedQuotaPackage(t *testing.T) {
	type args struct {
		premiumPackages []models.PremiumPackage
	}
	tests := []struct {
		name string
		u    *premiumPackageService
		args args
		want bool
	}{
		{
			name: "the package consists of an unlimited quota package",
			u:    &premiumPackageService{},
			args: args{
				premiumPackages: []models.PremiumPackage{
					{
						Type: models.UnilimitedQuota,
					},
					{
						Type: models.VerifiedUser,
					},
				},
			},
			want: true,
		},
		{
			name: "nil package",
			u:    &premiumPackageService{},
			args: args{
				premiumPackages: []models.PremiumPackage{},
			},
			want: false,
		},
		{
			name: "the package doesn't consists of an unlimited quota package",
			u:    &premiumPackageService{},
			args: args{
				premiumPackages: []models.PremiumPackage{
					{
						Type: models.VerifiedUser,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &premiumPackageService{}
			if got := u.IsConsistsUnlimitedQuotaPackage(tt.args.premiumPackages); got != tt.want {
				t.Errorf("premiumPackageService.IsConsistsUnlimitedQuotaPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_premiumPackageService_IsConsistsVerifiedUserPackage(t *testing.T) {
	type args struct {
		premiumPackages []models.PremiumPackage
	}
	tests := []struct {
		name string
		u    *premiumPackageService
		args args
		want bool
	}{
		{
			name: "the package consists of an verified user package",
			u:    &premiumPackageService{},
			args: args{
				premiumPackages: []models.PremiumPackage{
					{
						Type: models.UnilimitedQuota,
					},
					{
						Type: models.VerifiedUser,
					},
				},
			},
			want: true,
		},
		{
			name: "nil package",
			u:    &premiumPackageService{},
			args: args{
				premiumPackages: []models.PremiumPackage{},
			},
			want: false,
		},
		{
			name: "the package doesn't consists of an verified user package",
			u:    &premiumPackageService{},
			args: args{
				premiumPackages: []models.PremiumPackage{
					{
						Type: models.UnilimitedQuota,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &premiumPackageService{}
			if got := u.IsConsistsVerifiedUserPackage(tt.args.premiumPackages); got != tt.want {
				t.Errorf("premiumPackageService.IsConsistsVerifiedUserPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}
