package dto

import "github.com/google/uuid"

type GetListIdeasRequest struct {
	AuthorId *uuid.UUID `form:"author_id"`
	Search   *string    `form:"search"`
	StartAt  *uint      `form:"start_at"`
	Limit    *uint      `form:"limit"`
}

type CreateIdeaRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Content     *string `json:"content"`
	Category    *string `json:"category"`
}
