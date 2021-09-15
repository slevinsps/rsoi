package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Warranty struct {
	Comment      string    `json:"comment,omitempty"`
	ItemUID      uuid.UUID `json:"itemUid,omitempty"`
	WarrantyDate time.Time `json:"warrantyDate,omitempty"`
	Status       string    `json:"status,omitempty"`
}

type WarrantyParams struct {
	Reason         string `json:"reason,omitempty"`
	AvailableCount int    `json:"availableCount,omitempty"`
}

type WarrantyResponse struct {
	WarrantyDate string `json:"warrantyDate,omitempty"`
	Decision     string `json:"decision,omitempty"`
}
