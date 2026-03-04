package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	LeaderId    uuid.UUID `gorm:"column:leader_id;type:uuid; not null"`
	IdeaId      uuid.UUID `gorm:"column:idea_id;type:uuid"`
	Title       string    `gorm:"column:title; not null"`
	Description string    `gorm:"column:descriprion; not null"`
	Status      string    `gorm:"column:status; not null"`
	CreatedAt   time.Time `gorm:"column:created_at; not null"`

	Leader User `gorm:"foreignKey:LeaderId"`
	Idea   Idea `gorm:"foreignKey:IdeaId"`
}

func (p *Project) TableName() string { return "Project" }
