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

func (r *RepUsecase) CreateFile(file models.File) (err error) {
	return r.repo.CreateFile(file)
}

func (r *RepUsecase) UpdateFile(file models.File) (file_res models.File, err error) {
	return r.repo.UpdateFile(file)
}

func (r *RepUsecase) GetFileByFileUUID(fileUID uuid.UUID) (file models.File, checkFindFile bool, err error) {
	return r.repo.GetFileByFileUUID(fileUID)
}

func (r *RepUsecase) GetAllFiles() (files []models.File, err error) {
	return r.repo.GetAllFiles()
}

func (r *RepUsecase) GetAllFilesByEquipmentModelUUID(equipmentModelUUID uuid.UUID) (files []models.File, err error) {
	return r.repo.GetAllFilesByEquipmentModelUUID(equipmentModelUUID)
}

func (r *RepUsecase) DelFileByFileUUID(fileUID uuid.UUID) (err error) {
	return r.repo.DelFileByFileUUID(fileUID)
}
