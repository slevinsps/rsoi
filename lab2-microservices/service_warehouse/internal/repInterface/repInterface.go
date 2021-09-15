package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepInterface interface {
	InsertOrderItem(orderItemTableRecord models.OrderItemTableRecord) (err error)
	GetItemInfoByModeSize(model string, size string) (item models.Item, err error)
	GetItemInfoByOrderItemUID(orderItemUID uuid.UUID, availableCount bool) (item models.Item, err error)
	TakeOneItem(item models.Item) (err error)
	ReturnOneItem(orderItemUID uuid.UUID) (err error)
	CancelOrderItem(orderItemUID uuid.UUID) (err error)
}
