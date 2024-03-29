package profiles

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/github.com/vido21/dating-app/database"
	"github.com/github.com/vido21/dating-app/profiles/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func Test_profileService_GetProfileRecomendation(t *testing.T) {
	type args struct {
		excludeProfileIDs []uuid.UUID
		excludeUserID     uuid.UUID
	}
	tests := []struct {
		name    string
		u       *profileService
		args    args
		want    *models.Profile
		wantErr bool
		mock    func()
	}{
		{
			name: "success get data",
			u:    &profileService{},
			args: args{
				excludeProfileIDs: []uuid.UUID{},
				excludeUserID:     uuid.UUID{},
			},
			want: &models.Profile{},
			mock: func() {
				monkey.Patch(database.GetInstance, func() *gorm.DB {
					return &gorm.DB{}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Not", func(db *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
					return &gorm.DB{
						Error: nil,
					}
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
			u:    &profileService{},
			args: args{
				excludeProfileIDs: []uuid.UUID{},
				excludeUserID:     uuid.UUID{},
			},
			wantErr: true,
			mock: func() {
				monkey.Patch(database.GetInstance, func() *gorm.DB {
					return &gorm.DB{}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Not", func(db *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
					return &gorm.DB{
						Error: nil,
					}
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
			u := &profileService{}

			tt.mock()
			defer monkey.UnpatchAll()

			got, err := u.GetProfileRecomendation(tt.args.excludeProfileIDs, tt.args.excludeUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("profileService.GetProfileRecomendation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("profileService.GetProfileRecomendation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProfileService(t *testing.T) {
	tests := []struct {
		name string
		want ProfileService
		mock func()
	}{
		{
			name: "singleton is nil",
			want: &profileService{},
		},
		{
			name: "singleton is nil",
			want: &profileService{},
			mock: func() {
				singleton = &profileService{}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.mock != nil {
				tt.mock()
			}

			if got := GetProfileService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProfileService() = %v, want %v", got, tt.want)
			}
		})
	}
}
