package models

import (
	"github.com/github.com/vido21/dating-app/common/models"
)

const (
	UnilimitedQuota = "UNLIMITED_QUOTA"
	VerifiedUser    = "VERIFIED_USER"
)

type PremiumPackage struct {
	models.Base
	Description string  `json:"description"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
}
