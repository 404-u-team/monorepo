package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID       `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Email         string          `gorm:"column:email;unique;not null"`
	PasswordHash  string          `gorm:"column:password_hash;not null"`
	Nickname      string          `gorm:"column:nickname;unique;not null"`
	AvatarUrl     string          `gorm:"column:avatar_url;null"`
	MainRole      *uuid.UUID      `gorm:"column:main_role;null"`
	Bio           string          `gorm:"column:bio;null"`
	CreatedAt     time.Time       `gorm:"column:created_at;not null"`
	IsAdmin       bool            `gorm:"column:is_admin;not null"`
	MainRoleSkill *SkillCategory  `gorm:"foreignKey:MainRole;references:ID" json:"-"`
	Skills        []SkillCategory `gorm:"many2many:User_Skill;foreignKey:ID;joinForeignKey:user_id;references:ID;joinReferences:skill_id"`
}

func (u *User) TableName() string {
	return "User"
}
