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

func (r *RepUsecase) InsertOrderItem(orderItemTableRecord models.OrderItemTableRecord) (err error) {
	return r.repo.InsertOrderItem(orderItemTableRecord)
}

func (r *RepUsecase) GetItemInfoByModeSize(model string, size string) (item models.Item, err error) {
	return r.repo.GetItemInfoByModeSize(model, size)
}

func (r *RepUsecase) GetItemInfoByOrderItemUID(orderItemUID uuid.UUID, availableCount bool) (item models.Item, err error) {
	return r.repo.GetItemInfoByOrderItemUID(orderItemUID, availableCount)
}

func (r *RepUsecase) TakeOneItem(item models.Item) (err error) {
	return r.repo.TakeOneItem(item)
}

func (r *RepUsecase) ReturnOneItem(orderItemUID uuid.UUID) (err error) {
	return r.repo.ReturnOneItem(orderItemUID)
}
func (r *RepUsecase) CancelOrderItem(orderItemUID uuid.UUID) (err error) {
	return r.repo.CancelOrderItem(orderItemUID)
}
