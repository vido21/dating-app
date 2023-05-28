package models

import (
	"github.com/github.com/vido21/dating-app/common/models"
	uuid "github.com/satori/go.uuid"
)

type Blog struct {
	models.Base
	Title   string
	Content string
	UserID  uuid.UUID
}
