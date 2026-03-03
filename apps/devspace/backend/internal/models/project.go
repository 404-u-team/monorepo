package models

import (
	"time"
)

type Project struct {
	ID          uint      `gorm:"primaryKey; column:id"`
	LeaderId    uint      `gorm:"column:leader_id; not null"`
	IdeaId      uint      `gorm:"column:idea_id"`
	Title       string    `gorm:"column:title; not null"`
	Description string    `gorm:"column:descriprion; not null"`
	Status      string    `gorm:"column:status; not null"`
	CreatedAt   time.Time `gorm:"column:created_at; not null"`

	Leader User `gorm:"foreignKey:LeaderId"`
	Idea   Idea `gorm:"foreignKey:IdeaId"`
}

func (p *Project) TableName() string { return "Project" }
