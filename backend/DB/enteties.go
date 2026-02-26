package db

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey; column:id"`
	Email        string `gorm:"unique; column:email; not null"`
	PasswordHash string `gorm:"column:password_hash; not null"`
	Nickname     string `gorm:"column:nickname; not null"`
	AvatarUrl    string `gorm:"column:avatar_url; not null"`
	Status       string `gorm:"column:status; not null"`
	Bio          string `gorm:"column:bio; not null"`
}

type Notification struct {
	Id        uint      `gorm:"primaryKey; column:id"`
	UserId    uint      `gorm:"column:user_id; not null"`
	Message   string    `gorm:"column:message; not null"`
	IsRead    bool      `gorm:"column:is_read; not null"`
	CreatedAt time.Time `gorm:"column:created_at; not null"`

	User User `gorm:"foreignKey:UserId"`
}

type Idea struct {
	Id             uint   `gorm:"primaryKey; column:id"`
	AuthorId       uint   `gorm:"column:author_id"`
	Title          string `gorm:"column:title; not null"`
	Description    string `gorm:"column:description; not null"`
	ViewsCount     uint   `gorm:"column:views_id; not null"`
	FavoritesCount uint   `gorm:"column:favorites_count; not null"`

	Author User `gorm:"foreignKey:AuthorId"`
}

type Project struct {
	Id          uint   `gorm:"primaryKey; column:id"`
	LeaderId    uint   `gorm:"column:leader_id; not null"`
	IdeaId      uint   `gorm:"column:idea_id"`
	Title       string `gorm:"column:title; not null"`
	Description string `gorm:"column:descriprion"`
	Status      string `gorm:"column:status; not null"`
}
