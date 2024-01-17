package models

import (
	"github.com/github.com/vido21/dating-app/common/models"
	ProfileModels "github.com/github.com/vido21/dating-app/profiles/models"
	UserModels "github.com/github.com/vido21/dating-app/users/models"
	uuid "github.com/satori/go.uuid"
)

const (
	Like = "LIKE"
	Pass = "PASS"

	LimitSwipe int = 10
)

type Swipe struct {
	models.Base
	SwipeType string                  `json:"swipe_type"`
	UserID    uuid.UUID               `json:"user_id"`
	ProfileID uuid.UUID               `json:"profile_id"`
	Users     []UserModels.User       `json:"users"`
	Profile   []ProfileModels.Profile `json:"profiles"`
}
