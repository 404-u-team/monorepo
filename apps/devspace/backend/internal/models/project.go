package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	LeaderID    uuid.UUID  `gorm:"column:leader_id;type:uuid; not null"`
	IdeaID      *uuid.UUID `gorm:"column:idea_id;type:uuid"`
	Title       string     `gorm:"column:title;unique; not null"`
	Description *string    `gorm:"column:description"`
	Status      string     `gorm:"column:status; not null"`
	CreatedAt   time.Time  `gorm:"column:created_at; not null"`
	UpdatedAt   time.Time  `gorm:"column:updated_at; not null"`

	Leader User `gorm:"foreignKey:LeaderID"`
	Idea   Idea `gorm:"foreignKey:IdeaID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (p *Project) TableName() string { return "Project" }
