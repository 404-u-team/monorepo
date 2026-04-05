package dto

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
)

type GetIdeasRequest struct {
	AuthorId *uuid.UUID `form:"author_id"`
	Search   *string    `form:"search"`
	StartAt  *uint      `form:"start_at"`
	Limit    *uint      `form:"limit"`
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
	Total int64         `json:"total"`
	Ideas []models.Idea `json:"ideas"`
}
