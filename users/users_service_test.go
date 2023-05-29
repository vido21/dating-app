package users

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/github.com/vido21/dating-app/database"
	"github.com/github.com/vido21/dating-app/users/models"
	"github.com/jinzhu/gorm"
)

func Test_usersService_FindUserByEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		u    *usersService
		args args
		want *models.User
		mock func()
	}{
		{
			name: "success get data",
			u:    &usersService{},
			args: args{
				email: "user1@mail.com",
			},
			want: &models.User{},
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
			u:    &usersService{},
			args: args{
				email: "user1@mail.com",
			},
			want: nil,
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
			u := &usersService{}

			tt.mock()
			defer monkey.UnpatchAll()

			if got := u.FindUserByEmail(tt.args.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("usersService.FindUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_usersService_AddUser(t *testing.T) {
	type args struct {
		name     string
		email    string
		password string
	}
	tests := []struct {
		name string
		u    *usersService
		args args
		want *models.User
		mock func()
	}{
		{
			name: "success create data",
			u:    &usersService{},
			args: args{
				name:     "user1",
				email:    "user1@mail.com",
				password: "abcd",
			},
			want: &models.User{
				Name:     "user1",
				Email:    "user1@mail.com",
				Password: "abcd",
			},
			mock: func() {
				monkey.Patch(database.GetInstance, func() *gorm.DB {
					return &gorm.DB{}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Create", func(db *gorm.DB, value interface{}) *gorm.DB {
					return &gorm.DB{
						Error: nil,
					}
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &usersService{}

			tt.mock()
			defer monkey.UnpatchAll()

			if got := u.AddUser(tt.args.name, tt.args.email, tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("usersService.AddUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
