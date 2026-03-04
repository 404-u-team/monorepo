package models

import "github.com/google/uuid"

type SkillCategory struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ParentId uuid.UUID `gorm:"column:parent_id;type:uuid"`
	Name     string    `gorm:"column:name; not null"`

	Parent *SkillCategory `gorm:"foreignKey:ParentId"`
}

func (skillCategory *SkillCategory) TableName() string { return "Skill_Category" }
