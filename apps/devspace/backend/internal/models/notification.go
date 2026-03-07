package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid; not null"`
	Message   string    `gorm:"column:message; not null"`
	IsRead    bool      `gorm:"column:is_read; not null"`
	CreatedAt time.Time `gorm:"column:created_at; not null"`

	User User `gorm:"foreignKey:UserID"`
}

func (n *Notification) TableName() string { return "Notification" }
