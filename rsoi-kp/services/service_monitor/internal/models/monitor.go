package models

import (
	uuid "github.com/satori/go.uuid"
)

type Monitor struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	MonitorUUID uuid.UUID `json:"monitor_uuid,omitempty"`
	UserUUID    uuid.UUID `json:"user_uuid,omitempty"`
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
