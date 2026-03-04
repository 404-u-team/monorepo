package dto

import "github.com/google/uuid"

type CreateProjectRequest struct {
	IdeaID      *uuid.UUID `json:"idea_id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
}
