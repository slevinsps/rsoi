package models

import (
	uuid "github.com/satori/go.uuid"
)

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

type TokenDetails struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string    `json:"accessUuid,omitempty"`
	UserId     uuid.UUID `json:"userId,omitempty"`
}

type Service struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
