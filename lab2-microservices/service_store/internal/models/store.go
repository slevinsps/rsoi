package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Item struct {
	Model string `json:"model,omitempty"`
	Size  string `json:"size,omitempty"`
}

type OrderResponse struct {
	OrderUID       uuid.UUID `json:"orderUid,omitempty"`
	Date           string    `json:"date,omitempty"`
	Model          string    `json:"model,omitempty"`
	Size           string    `json:"size,omitempty"`
	WarrantyDate   string    `json:"warrantyDate,omitempty"`
	WarrantyStatus string    `json:"warrantyStatus,omitempty"`
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

type WarrantyResponse struct {
	OrderUID     uuid.UUID `json:"orderUid,omitempty"`
	WarrantyDate string    `json:"warrantyDate,omitempty"`
	Decision     string    `json:"decision,omitempty"`
}

type Warranty struct {
	Comment      string    `json:"comment,omitempty"`
	ItemUID      uuid.UUID `json:"itemUid,omitempty"`
	WarrantyDate time.Time `json:"warrantyDate,omitempty"`
	Status       string    `json:"status,omitempty"`
}

type User struct {
	ID       int       `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	UserUUID uuid.UUID `json:"user_uid,omitempty"`
}
