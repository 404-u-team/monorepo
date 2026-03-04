package models

type User struct {
	// 	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ID           uint   `gorm:"primaryKey; column:id"`
	Email        string `gorm:"unique; column:email; not null"`
	PasswordHash string `gorm:"column:password_hash; not null"`
	Nickname     string `gorm:"column:nickname; not null"`
	AvatarUrl    string `gorm:"column:avatar_url; null"`
	MainRole     string `gorm:"column:main_role; null"`
	Bio          string `gorm:"column:bio; null"`
}

func (u *User) TableName() string { return "User" }
