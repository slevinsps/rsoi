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

func (r *RepUsecase) GetDataByEquipmentUUID(equipmentUUID uuid.UUID) (data_arr []models.Data, err error) {
	return r.repo.GetDataByEquipmentUUID(equipmentUUID)
}

func (r *RepUsecase) DeleteDataByEquipmentUUID(equipmentUUID uuid.UUID) (err error) {
	return r.repo.DeleteDataByEquipmentUUID(equipmentUUID)
}

func (r *RepUsecase) CLearData() (err error) {
	return r.repo.CLearData()
}

func (r *RepUsecase) CreateData(data models.Data) (err error) {
	return r.repo.CreateData(data)
}
