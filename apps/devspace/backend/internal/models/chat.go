package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectId uuid.UUID `gorm:"column:project_id;type:uuid"`
	Title     string    `gorm:"column:title"`
	Type      string    `gorm:"column:type; not null"`
	CreatedAt time.Time `gorm:"column:created_at"`

	Project Project `gorm:"foreignKey:ProjectId"`
}

func (c *Chat) TableName() string { return "Chat" }
