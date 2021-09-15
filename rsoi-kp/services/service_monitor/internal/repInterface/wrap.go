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

func (r *RepUsecase) CreateMonitor(monitor models.Monitor) (err error) {
	return r.repo.CreateMonitor(monitor)
}

func (r *RepUsecase) GetMonitorByMonitorUUID(monitorUID uuid.UUID) (monitor models.Monitor, checkFindMonitor bool, err error) {
	return r.repo.GetMonitorByMonitorUUID(monitorUID)
}
func (r *RepUsecase) GetMonitorByMonitorUUIDuserUUID(monitorUID uuid.UUID, userUID uuid.UUID) (monitor models.Monitor, checkFindMonitor bool, err error) {
	return r.repo.GetMonitorByMonitorUUIDuserUUID(monitorUID, userUID)
}
func (r *RepUsecase) GetAllMonitorsByUserUUID(userUID uuid.UUID) (monitors []models.Monitor, err error) {
	return r.repo.GetAllMonitorsByUserUUID(userUID)
}

func (r *RepUsecase) DelMonitorByMonitorUUID(monitorUID uuid.UUID) (err error) {
	return r.repo.DelMonitorByMonitorUUID(monitorUID)
}

func (r *RepUsecase) AddEquipment(monitorUUID uuid.UUID, equipmentUUID uuid.UUID) (err error) {
	return r.repo.AddEquipment(monitorUUID, equipmentUUID)
}

func (r *RepUsecase) DelEquipment(monitorUUID uuid.UUID, equipmentUUID uuid.UUID) (err error) {
	return r.repo.DelEquipment(monitorUUID, equipmentUUID)
}
