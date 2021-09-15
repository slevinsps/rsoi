package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepInterface interface {
	OrderCreate(orderUID uuid.UUID, userUID uuid.UUID, itemUID uuid.UUID) (err error)
	GetUserOrderInfo(userUID uuid.UUID, orderUID uuid.UUID) (order models.Order, err error)
	GetUserOrdersInfo(userUID uuid.UUID) (orders []models.Order, err error)
	GetOrderInfoByOrderUID(orderUID uuid.UUID) (order models.Order, err error)
}
