package repInterface

import "restService/internal/models"

type RepInterface interface {
	PersonCreate(person models.Person) (models.Person, error)
	GetPersonByID(id int) (models.Person, bool, error)
	GetAllPersonsInfo() ([]models.Person, error)
	UpdatePersonInfo(person models.Person) (models.Person, error)
	DeletePersonInfo(id int) error
}
