package models

import "github.com/google/uuid"

type ProjectSlot struct {
	ID              uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID       uuid.UUID  `gorm:"column:project_id;type:uuid;not null;uniqueIndex:idx_project_user,where:user_id IS NOT NULL"`
	SkillCategoryID uuid.UUID  `gorm:"column:skill_category_id;type:uuid; not null"`
	UserID          *uuid.UUID `gorm:"column:project_id;type:uuid;not null;uniqueIndex:idx_project_user,where:user_id IS NOT NULL"`
	Title           string     `gorm:"column:title; not null"`
	Description     *string    `gorm:"column:description"`
	Status          string     `gorm:"column:status; not null"`

	Project Project       `gorm:"foreignKey:ProjectID"`
	Skill   SkillCategory `gorm:"foreignKey:SkillCategoryID"`
	User    User          `gorm:"foreignKey:UserID"`
}

func (ps *ProjectSlot) TableName() string { return "Project_Slot" }
