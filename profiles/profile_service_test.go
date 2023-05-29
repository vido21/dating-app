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
		excludeUserIDs []uuid.UUID
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
				excludeUserIDs: []uuid.UUID{},
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
				excludeUserIDs: []uuid.UUID{},
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

			got, err := u.GetProfileRecomendation(tt.args.excludeUserIDs)
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
