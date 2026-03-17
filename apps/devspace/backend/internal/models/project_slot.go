package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectSlot struct {
	ID              uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProjectID       uuid.UUID  `gorm:"column:project_id;type:uuid;not null;uniqueIndex:idx_project_user,where:project_id IS NOT NULL" json:"project_id"`
	SkillCategoryID uuid.UUID  `gorm:"column:skill_category_id;type:uuid; not null" json:"skill_category_id"`
	UserID          *uuid.UUID `gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_project_user,where:user_id IS NOT NULL" json:"user_id"`
	Title           string     `gorm:"column:title; not null" json:"title"`
	Description     *string    `gorm:"column:description" json:"description"`
	Status          string     `gorm:"column:status; not null" json:"status"`
	CreatedAt       time.Time  `gorm:"column:created_at; not null" json:"created_at"`

	Project Project       `gorm:"foreignKey:ProjectID" json:"-"`
	Skill   SkillCategory `gorm:"foreignKey:SkillCategoryID" json:"-"`
	User    User          `gorm:"foreignKey:UserID" json:"-"`
}

func (ps *ProjectSlot) TableName() string { return "Project_Slot" }
