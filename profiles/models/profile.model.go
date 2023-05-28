package models

import (
	"github.com/github.com/vido21/dating-app/common/models"
	uuid "github.com/satori/go.uuid"
)

// Sex Type
const (
	FemaleIntType int = 0
	MaleIntType   int = 1

	MaleSexTypeString    string = "MALE"
	FemaleSexTypeString  string = "FEMALE"
	DefaultSexTypeString string = ""
)

var (
	Sex = map[string]int{
		FemaleSexTypeString:  FemaleIntType,
		MaleSexTypeString:    MaleIntType,
		DefaultSexTypeString: FemaleIntType,
	}
)

type Profile struct {
	models.Base
	ProfilePicture string    `json:"profile_picture"`
	Sex            int       `json:"sex"`
	About          string    `json:"about"`
	UserID         uuid.UUID `json:"user_id,omitempty"`
}
