package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetIdeasRequest struct {
	AuthorId   *QueryUUID `form:"author_id"`
	Search     *string    `form:"search"`
	StartAt    *uint      `form:"start_at"`
	Limit      *uint      `form:"limit"`
	Views      *string    `form:"views" binding:"omitempty,oneof=asc desc"`
	Favorites  *string    `form:"favorites" binding:"omitempty,oneof=asc desc"`
	IsFavorite bool       `form:"is_favorite"`
}

type IdeaBlock struct {
	ID             uuid.UUID `json:"id"`
	AuthorID       uuid.UUID `json:"author_id"`
	IsAuthor       bool      `json:"is_author"`
	IsFavorite     bool      `json:"is_favorite"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ViewsCount     uint      `json:"views_count"`
	FavoritesCount uint      `json:"favorites_count"`
	Category       string    `json:"category"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type GetIdeaResponse struct {
	ID             uuid.UUID `json:"id"`
	AuthorID       uuid.UUID `json:"author_id"`
	IsAuthor       bool      `json:"is_author"`
	IsFavorite     bool      `json:"is_favorite"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Content        *string   `json:"content"`
	ViewsCount     uint      `json:"views_count"`
	FavoritesCount uint      `json:"favorites_count"`
	Category       string    `json:"category"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateIdeaRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Content     *string `json:"content"`
	Category    *string `json:"category"`
}

type UpdateIdeaRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=3,max=255"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	Content     *string `json:"content" binding:"omitempty,max=10000"`
}

type GetIdeasResponse struct {
	Total int64       `json:"total"`
	Ideas []IdeaBlock `json:"ideas"`
}

type ToggleFavoriteResponse struct {
	IsFavorite bool `json:"is_favorite"`
}
