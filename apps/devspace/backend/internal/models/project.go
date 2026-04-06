package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LeaderID    uuid.UUID  `gorm:"column:leader_id;type:uuid; not null" json:"leader_id"`
	IdeaID      *uuid.UUID `gorm:"column:idea_id;type:uuid" json:"idea_id"`
	Title       string     `gorm:"column:title;unique; not null" json:"title"`
	Description *string    `gorm:"column:description" json:"description"`
	Content     *string    `gorm:"column:content" json:"content"`
	Status      string     `gorm:"column:status; not null" json:"status"`
	CreatedAt   time.Time  `gorm:"column:created_at; not null" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at; not null" json:"updated_at"`

	Leader       User          `gorm:"foreignKey:LeaderID" json:"-"`
	Idea         Idea          `gorm:"foreignKey:IdeaID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ProjectSlots []ProjectSlot `gorm:"foreignKey:ProjectID" json:"-"`
}

func (p *Project) TableName() string { return "Project" }
