package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepUsecase struct {
	repo RepInterface
}

func NewPUsecase(repo RepInterface) *RepUsecase {
	return &RepUsecase{repo: repo}
}

func (r *RepUsecase) GetUserByUUID(userUID uuid.UUID) (user models.User, err error) {
	return r.repo.GetUserByUUID(userUID)
}
