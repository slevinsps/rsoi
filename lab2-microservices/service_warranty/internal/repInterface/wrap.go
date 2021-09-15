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

func (r *RepUsecase) StartWarrantyPeriod(itemUID uuid.UUID) (err error) {
	return r.repo.StartWarrantyPeriod(itemUID)
}

func (r *RepUsecase) GetWarranty(itemUID uuid.UUID) (warranty models.Warranty, err error) {
	return r.repo.GetWarranty(itemUID)
}

func (r *RepUsecase) UpdateWarranty(warranty models.Warranty) (err error) {
	return r.repo.UpdateWarranty(warranty)
}

func (r *RepUsecase) CloseWarranty(itemUID uuid.UUID) (err error) {
	return r.repo.CloseWarranty(itemUID)
}
