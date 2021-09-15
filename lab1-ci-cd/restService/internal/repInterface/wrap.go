package repInterface

import (
	"restService/internal/models"
)

type RepUsecase struct {
	repo RepInterface
}

func NewPUsecase(repo RepInterface) *RepUsecase {
	return &RepUsecase{repo: repo}
}

func (r *RepUsecase) PersonCreate(person models.Person) (models.Person, error) {
	return r.repo.PersonCreate(person)
}

func (r *RepUsecase) GetPersonByID(id int) (models.Person, bool, error) {
	return r.repo.GetPersonByID(id)
}

func (r *RepUsecase) GetAllPersonsInfo() ([]models.Person, error) {
	return r.repo.GetAllPersonsInfo()
}

func (r *RepUsecase) UpdatePersonInfo(person models.Person) (models.Person, error) {
	return r.repo.UpdatePersonInfo(person)
}

func (r *RepUsecase) DeletePersonInfo(id int) error {
	return r.repo.DeletePersonInfo(id)
}
