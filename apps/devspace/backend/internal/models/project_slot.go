package models

import "github.com/google/uuid"

type ProjectSlot struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectId       uuid.UUID `gorm:"column:project_id;type:uuid; not null"`
	SkillCategoryId uuid.UUID `gorm:"column:skill_category_id;type:uuid; not null"`
	UserId          uuid.UUID `gorm:"column:user_id;type:uuid"`
	Status          string    `gorm:"column:status; not null"`

	Project Project       `gorm:"foreignKey:ProjectId"`
	Skill   SkillCategory `gorm:"foreignKey:SkillCategoryId"`
	User    User          `gorm:"foreignKey:UserId"`
}

func (ps *ProjectSlot) TableName() string { return "Project_Slot" }
