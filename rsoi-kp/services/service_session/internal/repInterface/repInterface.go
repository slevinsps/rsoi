package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepInterface interface {
	GetUserByUUID(userUID uuid.UUID) (user models.User, checkFindUser bool, err error)
	GetUserByLogin(login string) (user models.User, checkFindUser bool, err error)
	CreateUser(user models.User) (checkFindUser bool, err error)
	GetAllUsers() (users []models.User, err error)
}
