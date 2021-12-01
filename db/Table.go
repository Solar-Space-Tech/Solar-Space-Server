package db

import (
	"time"

	"github.com/fox-one/pkg/property"
)

type User struct {
	Uuid  string `gorm:"primaryKey"`
	Phone string `gorm:"unique"`
	Name  string
}

type SubWallet struct {
	ClientID   string `json:"client_id"`
	SessionID  string `json:"session_id"`
	PrivateKey string `gorm:"type:TEXT"`
	PinToken   string `json:"pin_token"`
}

type Property struct {
	Key       string         `gorm:"size:64;PRIMARY_KEY"`
	Value     property.Value `gorm:"type:varchar(256)"`
	UpdatedAt time.Time      `gorm:"precision:6"`
}
