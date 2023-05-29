package auth

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/dgrijalva/jwt-go"

	UserModels "github.com/github.com/vido21/dating-app/users/models"
)

func Test_authService_GetAccessToken(t *testing.T) {
	type args struct {
		user *UserModels.User
	}
	tests := []struct {
		name    string
		s       *authService
		args    args
		want    string
		wantErr bool
		mock    func()
	}{
		{
			name: "default case",
			s:    &authService{},
			args: args{
				user: &UserModels.User{},
			},
			want: "token",
			mock: func() {
				monkey.Patch(jwt.NewWithClaims, func(signingMethod jwt.SigningMethod, claims jwt.Claims) *jwt.Token {
					return &jwt.Token{}
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(&jwt.Token{}), "SignedString", func(token *jwt.Token, key interface{}) (string, error) {
					return "token", nil
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authService{}

			tt.mock()
			defer monkey.UnpatchAll()

			got, err := s.GetAccessToken(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("authService.GetAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
