package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepInterface interface {
	CreateFile(user models.File) (err error)
	UpdateFile(user models.File) (file_res models.File, err error)
	GetAllFiles() (files []models.File, err error)
	GetAllFilesByEquipmentModelUUID(equipmentModelUUID uuid.UUID) (files []models.File, err error)
	GetFileByFileUUID(fileUID uuid.UUID) (file models.File, checkFindFile bool, err error)
	DelFileByFileUUID(fileUID uuid.UUID) (err error)
}
