package models

import "github.com/google/uuid"

type Idea struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	AuthorID       uuid.UUID `gorm:"column:author_id;type:uuid"`
	Title          string    `gorm:"column:title; not null; unique"`
	Description    string    `gorm:"column:description; not null"`
	ViewsCount     uint      `gorm:"column:views_count; not null"`
	FavoritesCount uint      `gorm:"column:favorites_count; not null"`
	Category       string    `gorm:"column:category; not null"`

	Author User `gorm:"foreignKey:AuthorID"`
}

func (i *Idea) TableName() string { return "Idea" }
