package models

import "github.com/google/uuid"

type SkillCategory struct {
	ID       uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ParentID *uuid.UUID     `gorm:"column:parent_id;type:uuid" json:"parent_id"`
	Name     string         `gorm:"column:name; not null" json:"name"`
	Icon     string         `gorm:"column:icon;type:text" json:"icon"`
	Color    string         `gorm:"column:color;type:text" json:"color"`
	Parent   *SkillCategory `gorm:"foreignKey:ParentID" json:"-"`
}

func (skillCategory *SkillCategory) TableName() string { return "Skill_Category" }
