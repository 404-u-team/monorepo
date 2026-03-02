package models

import (
	"time"
)

type Chat struct {
	Id        uint      `gorm:"primaryKey; column:id"`
	ProjectId uint      `gorm:"column:project_id"`
	Title     string    `gorm:"column:title"`
	Type      string    `gorm:"column:type; not null"`
	CreatedAt time.Time `gorm:"column:created_at"`

	Project Project `gorm:"foreignKey:ProjectId"`
}

func (c *Chat) TableName() string { return "Chat" }
