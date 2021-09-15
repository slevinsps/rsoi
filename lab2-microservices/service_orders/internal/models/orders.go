package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Item struct {
	Model string `json:"model,omitempty"`
	Size  string `json:"size,omitempty"`
}

type OrderAll struct {
	ID        int       `json:"id,omitempty"`
	ItemUID   uuid.UUID `json:"itemUid,omitempty"`
	OrderDate time.Time `json:"orderDate,omitempty"`
	OrderUID  uuid.UUID `json:"orderUid,omitempty"`
	Status    string    `json:"status,omitempty"`
	UserUID   uuid.UUID `json:"userUid,omitempty"`
}

type Order struct {
	OrderUID  uuid.UUID `json:"orderUid,omitempty"`
	OrderDate time.Time `json:"orderDate,omitempty"`
	ItemUID   uuid.UUID `json:"itemUid,omitempty"`
	Status    string    `json:"status,omitempty"`
}

type ItemOrder struct {
	OrderItemUid uuid.UUID `json:"orderItemUid,omitempty"`
	OrderUid     uuid.UUID `json:"orderUid,omitempty"`
	Model        string    `json:"model,omitempty"`
	Size         string    `json:"size,omitempty"`
}

type OrderUid struct {
	OrderUid uuid.UUID `json:"orderUid"`
}

type WarrantyParams struct {
	Reason         string `json:"reason,omitempty"`
	AvailableCount int    `json:"availableCount,omitempty"`
}
