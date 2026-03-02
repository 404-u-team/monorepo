package models

type User struct {
	ID           uint   `gorm:"primaryKey; column:id"`
	Email        string `gorm:"unique; column:email; not null"`
	PasswordHash string `gorm:"column:password_hash; not null"`
	Nickname     string `gorm:"column:nickname; not null"`
	AvatarUrl    string `gorm:"column:avatar_url; not null"`
	Status       string `gorm:"column:status; not null"`
	Bio          string `gorm:"column:bio; not null"`
}

func (u *User) TableName() string { return "User" }
