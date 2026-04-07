package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetIdeasRequest struct {
	AuthorId   *uuid.UUID `form:"author_id"`
	Search     *string    `form:"search"`
	StartAt    *uint      `form:"start_at"`
	Limit      *uint      `form:"limit"`
	Views      *string    `form:"views" binding:"omitempty,oneof=asc desc"`
	Favorites  *string    `form:"favorites" binding:"omitempty,oneof=asc desc"`
	IsFavorite bool       `form:"is_favorite"`
}

type IdeaBlock struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	IsAuthor       bool      `gorm:"column:is_author" json:"is_author"`
	IsFavorite     bool      `gorm:"column:is_favorite" json:"is_favorite"`
	Title          string    `gorm:"column:title; not null; unique" json:"title"`
	Description    string    `gorm:"column:description; not null" json:"description"`
	ViewsCount     uint      `gorm:"column:views_count; not null" json:"views_count"`
	FavoritesCount uint      `gorm:"column:favorites_count; not null" json:"favorites_count"`
	Category       string    `gorm:"column:category; not null" json:"category"`
	CreatedAt      time.Time `gorm:"column:created_at; not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at; not null" json:"updated_at"`
}

type CreateIdeaRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Content     *string `json:"content"`
	Category    *string `json:"category"`
}

type UpdateIdeaRequest struct {
	Title       *string `json:"title" binding:"min=3,max=255"`
	Description *string `json:"description" binding:"min=3,max=255"`
}

type GetIdeasResponse struct {
	Total int64       `json:"total"`
	Ideas []IdeaBlock `json:"ideas"`
}

type ToggleFavoriteResponse struct {
	IsFavorite bool `json:"is_favorite"`
}
