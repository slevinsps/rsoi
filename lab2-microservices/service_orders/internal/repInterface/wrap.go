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

func (r *RepUsecase) OrderCreate(orderUID uuid.UUID, userUID uuid.UUID, itemUID uuid.UUID) (err error) {
	return r.repo.OrderCreate(orderUID, userUID, itemUID)
}

func (r *RepUsecase) GetUserOrdersInfo(userUID uuid.UUID) (orders []models.Order, err error) {
	return r.repo.GetUserOrdersInfo(userUID)
}

func (r *RepUsecase) GetUserOrderInfo(userUID uuid.UUID, orderUID uuid.UUID) (order models.Order, err error) {
	return r.repo.GetUserOrderInfo(userUID, orderUID)
}

func (r *RepUsecase) GetOrderInfoByOrderUID(orderUID uuid.UUID) (order models.Order, err error) {
	return r.repo.GetOrderInfoByOrderUID(orderUID)
}

// func (r *RepUsecase) GetPersonByID(id int) (models.Person, bool, error) {
// 	return r.repo.GetPersonByID(id)
// }

// func (r *RepUsecase) GetAllPersonsInfo() ([]models.Person, error) {
// 	return r.repo.GetAllPersonsInfo()
// }

// func (r *RepUsecase) UpdatePersonInfo(person models.Person) (models.Person, error) {
// 	return r.repo.UpdatePersonInfo(person)
// }

// func (r *RepUsecase) DeletePersonInfo(id int) error {
// 	return r.repo.DeletePersonInfo(id)
// }
