package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepInterface interface {
	GetDataByEquipmentUUID(equipmentUUID uuid.UUID) (data_arr []models.Data, err error)
	DeleteDataByEquipmentUUID(equipmentUUID uuid.UUID) (err error)
	CLearData() (err error)
	CreateData(data models.Data) (err error)
}
