package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepInterface interface {
	GetUserByUUID(userUID uuid.UUID) (user models.User, err error)
}
