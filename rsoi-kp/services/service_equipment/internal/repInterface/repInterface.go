package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepInterface interface {
	CreateEquipment(user models.Equipment) (err error)
	UpdateEquipment(user models.Equipment) (equipment_res models.Equipment, err error)
	GetAllEquipments() (equipments []models.Equipment, err error)
	GetEquipmentByEquipmentUUID(equipmentUID uuid.UUID) (equipment models.Equipment, checkFindEquipment bool, err error)
	DelEquipmentByEquipmentUUID(equipmentUID uuid.UUID) (err error)
	GetEquipmentsByMonitorUUID(monitorUUID uuid.UUID) (equipments []models.Equipment, err error)
	GetNotAddedEquipmentsByMonitorUUID(monitorUUID uuid.UUID) (equipments []models.Equipment, err error)

	CreateEquipmentModel(user models.EquipmentModel) (err error)
	GetAllEquipmentModels() (equipments []models.EquipmentModel, err error)
	GetEquipmentModelByEquipmentModelUUID(equipmentModelUID uuid.UUID) (equipmentModel models.EquipmentModel, checkFindEquipmentModel bool, err error)
	DelEquipmentModelByEquipmentModelUUID(equipmentModelUID uuid.UUID) (err error)
}
