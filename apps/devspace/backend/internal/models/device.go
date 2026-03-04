package models

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserId       uuid.UUID `gorm:"column:user_id;type:uuid; not null"`
	DeviceName   string    `gorm:"column:device_name; not null"`
	Ip           string    `gorm:"column:ip; not null"`
	RefreshToken string    `gorm:"column:refresh_token; not null"`
	CreatedAt    time.Time `gorm:"column:created_at; not null"`
	ValidUntil   time.Time `gorm:"column:valid_until; not null"`
}

func (d *Device) TableName() string { return "Device" }
