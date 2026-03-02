package models

import (
	"time"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey; column:id"`
	UserId    uint      `gorm:"column:user_id; not null"`
	Message   string    `gorm:"column:message; not null"`
	IsRead    bool      `gorm:"column:is_read; not null"`
	CreatedAt time.Time `gorm:"column:created_at; not null"`

	User User `gorm:"foreignKey:UserId"`
}

func (n *Notification) TableName() string { return "Notification" }
