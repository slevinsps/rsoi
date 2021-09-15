package models

import (
	uuid "github.com/satori/go.uuid"
)

type TokenDetails struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	AccessUuid   string `json:"access_uuid,omitempty"`
	RefreshUuid  string `json:"refresh_uuid,omitempty"`
	AtExpires    int64  `json:"at_expires,omitempty"`
	RtExpires    int64  `json:"rt_expires,omitempty"`
}

type AccessDetails struct {
	AccessUuid string    `json:"accessUuid,omitempty"`
	UserId     uuid.UUID `json:"userId,omitempty"`
}

type Monitor struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	MonitorUUID uuid.UUID `json:"monitor_uuid,omitempty"`
	UserUUID    uuid.UUID `json:"user_uuid,omitempty"`
}

type Equipment struct {
	ID                 int       `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	EquipmentUUID      uuid.UUID `json:"equipment_uuid,omitempty"`
	EquipmentModelUUID uuid.UUID `json:"equipment_model_uuid,omitempty"`
	Status             string    `json:"status,omitempty"`
}

type EquipmentModel struct {
	ID                 int       `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	EquipmentModelUUID uuid.UUID `json:"equipment_model_uuid,omitempty"`
	ReleaseDate        string    `json:"release_date,omitempty"`
}

type File struct {
	ID                 int       `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	FileUUID           uuid.UUID `json:"file_uuid,omitempty"`
	Path               string    `json:"path,omitempty"`
	EquipmentModelUUID uuid.UUID `json:"equipment_model_uuid,omitempty"`
}

type FileSend struct {
	Name               string    `json:"name,omitempty"`
	FileUUID           uuid.UUID `json:"file_uuid,omitempty"`
	EquipmentModelUUID uuid.UUID `json:"equipment_model_uuid,omitempty"`
}

type Data struct {
	ID            int       `json:"id,omitempty"`
	DataUUID      uuid.UUID `json:"data_uuid,omitempty"`
	EquipmentUUID uuid.UUID `json:"equipment_uuid,omitempty"`
	Temperature   float32   `json:"temperature,omitempty"`
	Voltage       float32   `json:"voltage,omitempty"`
	Frequency     float32   `json:"frequency,omitempty"`
	LoadLevel     float32   `json:"load_level,omitempty"`
}

type Service struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
