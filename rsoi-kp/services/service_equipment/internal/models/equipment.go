package models

import (
	uuid "github.com/satori/go.uuid"
)

type Equipment struct {
	ID                 int       `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	EquipmentUUID      uuid.UUID `json:"equipment_uuid,omitempty"`
	ModelName          string    `json:"model_name,omitempty"`
	EquipmentModelUUID uuid.UUID `json:"equipment_model_uuid,omitempty"`
	Status             string    `json:"status,omitempty"`
}

type EquipmentModel struct {
	ID                 int       `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	EquipmentModelUUID uuid.UUID `json:"equipment_model_uuid,omitempty"`
}

type Service struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
