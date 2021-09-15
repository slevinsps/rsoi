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

func (r *RepUsecase) CreateEquipment(equipment models.Equipment) (err error) {
	return r.repo.CreateEquipment(equipment)
}

func (r *RepUsecase) UpdateEquipment(equipment models.Equipment) (equipment_res models.Equipment, err error) {
	return r.repo.UpdateEquipment(equipment)
}

func (r *RepUsecase) GetEquipmentByEquipmentUUID(equipmentUID uuid.UUID) (equipment models.Equipment, checkFindEquipment bool, err error) {
	return r.repo.GetEquipmentByEquipmentUUID(equipmentUID)
}

func (r *RepUsecase) GetAllEquipments() (equipments []models.Equipment, err error) {
	return r.repo.GetAllEquipments()
}

func (r *RepUsecase) DelEquipmentByEquipmentUUID(equipmentUID uuid.UUID) (err error) {
	return r.repo.DelEquipmentByEquipmentUUID(equipmentUID)
}

func (r *RepUsecase) CreateEquipmentModel(equipmentModel models.EquipmentModel) (err error) {
	return r.repo.CreateEquipmentModel(equipmentModel)
}

func (r *RepUsecase) GetEquipmentModelByEquipmentModelUUID(equipmentModelUID uuid.UUID) (equipmentModel models.EquipmentModel, checkFindEquipmentModel bool, err error) {
	return r.repo.GetEquipmentModelByEquipmentModelUUID(equipmentModelUID)
}

func (r *RepUsecase) GetAllEquipmentModels() (equipmentModels []models.EquipmentModel, err error) {
	return r.repo.GetAllEquipmentModels()
}

func (r *RepUsecase) DelEquipmentModelByEquipmentModelUUID(equipmentModelUID uuid.UUID) (err error) {
	return r.repo.DelEquipmentModelByEquipmentModelUUID(equipmentModelUID)
}

func (r *RepUsecase) GetEquipmentsByMonitorUUID(monitorUUID uuid.UUID) (equipments []models.Equipment, err error) {
	return r.repo.GetEquipmentsByMonitorUUID(monitorUUID)
}

func (r *RepUsecase) GetNotAddedEquipmentsByMonitorUUID(monitorUUID uuid.UUID) (equipments []models.Equipment, err error) {
	return r.repo.GetNotAddedEquipmentsByMonitorUUID(monitorUUID)
}
