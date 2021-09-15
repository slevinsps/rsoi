package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepInterface interface {
	CreateMonitor(user models.Monitor) (err error)
	GetAllMonitorsByUserUUID(userUID uuid.UUID) (monitors []models.Monitor, err error)
	GetMonitorByMonitorUUID(monitorUID uuid.UUID) (monitor models.Monitor, checkFindMonitor bool, err error)
	DelMonitorByMonitorUUID(monitorUID uuid.UUID) (err error)
	AddEquipment(monitorUUID uuid.UUID, equipmentUUID uuid.UUID) (err error)
	DelEquipment(monitorUUID uuid.UUID, equipmentUUID uuid.UUID) (err error)
	GetMonitorByMonitorUUIDuserUUID(monitorUID uuid.UUID, userUID uuid.UUID) (monitor models.Monitor, checkFindMonitor bool, err error)
}
