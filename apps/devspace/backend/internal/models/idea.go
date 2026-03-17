package models

import (
	"time"

	"github.com/google/uuid"
)

type Idea struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AuthorID       uuid.UUID `gorm:"column:author_id;type:uuid" json:"author_id"`
	Title          string    `gorm:"column:title; not null; unique" json:"title"`
	Description    string    `gorm:"column:description; not null" json:"description"`
	Content        *string   `gorm:"column:content" json:"content"`
	ViewsCount     uint      `gorm:"column:views_count; not null" json:"views_count"`
	FavoritesCount uint      `gorm:"column:favorites_count; not null" json:"favorites_count"`
	Category       string    `gorm:"column:category; not null" json:"category"`
	CreatedAt      time.Time `gorm:"column:created_at; not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at; not null" json:"updated_at"`

	Author User `gorm:"foreignKey:AuthorID" json:"-"`
}

func (i *Idea) TableName() string { return "Idea" }
