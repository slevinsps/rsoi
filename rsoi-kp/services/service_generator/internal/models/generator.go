package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Data struct {
	ID            int       `json:"id,omitempty"`
	DataUUID      uuid.UUID `json:"data_uuid,omitempty"`
	EquipmentUUID uuid.UUID `json:"equipment_uuid,omitempty"`
	Temperature   float32   `json:"temperature,omitempty"`
	Voltage       float32   `json:"voltage,omitempty"`
	Frequency     float32   `json:"frequency,omitempty"`
	LoadLevel     float32   `json:"load_level,omitempty"`
	Timestamp     time.Time `json:"timestamp,omitempty"`
}

type Equipment struct {
	ID                 int       `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	EquipmentUUID      uuid.UUID `json:"equipment_uuid,omitempty"`
	EquipmentModelUUID uuid.UUID `json:"equipment_model_uuid,omitempty"`
	Status             string    `json:"status,omitempty"`
}

type Service struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type TokenDetails struct {
	AccessToken string `json:"access_token,omitempty"`
}
