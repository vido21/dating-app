package profiles

import (
	"sync"

	"github.com/github.com/vido21/dating-app/database"
	"github.com/github.com/vido21/dating-app/profiles/models"
	uuid "github.com/satori/go.uuid"
)

type profileService struct{}

var singleton ProfileService
var once sync.Once

func GetProfileService() ProfileService {
	if singleton != nil {
		return singleton
	}
	once.Do(func() {
		singleton = &profileService{}
	})
	return singleton
}

func SetProfileService(service ProfileService) ProfileService {
	original := singleton
	singleton = service
	return original
}

type ProfileService interface {
	GetProfileRecomendation(excludeUserIDs []uuid.UUID) (*models.Profile, error)
}

func (u *profileService) GetProfileRecomendation(excludeUserIDs []uuid.UUID) (*models.Profile, error) {
	db := database.GetInstance()
	var profile models.Profile

	err := db.Not("user_id IN (?)", excludeUserIDs).First(&profile).Error
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
