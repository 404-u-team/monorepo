package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"unique; column:email; not null"`
	PasswordHash string    `gorm:"column:password_hash; not null"`
	Nickname     string    `gorm:"column:nickname; not null"`
	AvatarUrl    string    `gorm:"column:avatar_url; null"`
	MainRole     string    `gorm:"column:main_role; null"`
	Bio          string    `gorm:"column:bio; null"`
	IsAdmin      bool      `gorm:"column:is_admin; not null"`
}

func (u *User) TableName() string { return "User" }
