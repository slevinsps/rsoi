package models

import (
	uuid "github.com/satori/go.uuid"
)

type OrderItem struct {
	OrderItemUid uuid.UUID `json:"orderItemUid,omitempty"`
	OrderUid     uuid.UUID `json:"orderUid,omitempty"`
	Model        string    `json:"model,omitempty"`
	Size         string    `json:"size,omitempty"`
}

type OrderItemTableRecord struct {
	Canceled     bool      `json:"canceled,omitempty"`
	OrderItemUid uuid.UUID `json:"orderItemUid,omitempty"`
	OrderUid     uuid.UUID `json:"orderUid,omitempty"`
	ItemID       int       `json:"itemId,omitempty"`
}

type Item struct {
	ID             int    `json:"item_id,omitempty"`
	Model          string `json:"model,omitempty"`
	Size           string `json:"size,omitempty"`
	AvailableCount int    `json:"availableCount,omitempty"`
}

type WarrantyParams struct {
	Reason         string `json:"reason,omitempty"`
	AvailableCount int    `json:"availableCount,omitempty"`
}
