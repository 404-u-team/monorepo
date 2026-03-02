package models

import (
	"time"
)

type Device struct {
	Id           uint      `gorm:"primaryKey; column:id"`
	UserId       uint      `gorm:"column:user_id; not null"`
	DeviceName   string    `gorm:"column:device_name; not null"`
	Ip           string    `gorm:"column:ip; not null"`
	RefreshToken string    `gorm:"column:refresh_token; not null"`
	CreatedAt    time.Time `gorm:"column:created_at; not null"`
	ValidUntil   time.Time `gorm:"column:valid_until; not null"`
}

func (d *Device) TableName() string { return "Device" }
