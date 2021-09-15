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

func (r *RepUsecase) GetUserByUUID(userUID uuid.UUID) (user models.User, checkFindUser bool, err error) {
	return r.repo.GetUserByUUID(userUID)
}

func (r *RepUsecase) GetUserByLogin(login string) (user models.User, checkFindUser bool, err error) {
	return r.repo.GetUserByLogin(login)
}

func (r *RepUsecase) CreateUser(user models.User) (checkFindUser bool, err error) {
	return r.repo.CreateUser(user)
}

func (r *RepUsecase) GetAllUsers() (users []models.User, err error) {
	return r.repo.GetAllUsers()
}
