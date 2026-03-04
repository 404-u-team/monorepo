package models

import "github.com/google/uuid"

type UserSkill struct {
	UserID  uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey"`
	SkillId uuid.UUID `gorm:"column:skill_id;type:uuid;primaryKey"`

	User  User          `gorm:"foreignKey:UserId"`
	Skill SkillCategory `gorm:"foreignKey:SkillId"`
}

func (us *UserSkill) TableName() string { return "User_Skill" }
