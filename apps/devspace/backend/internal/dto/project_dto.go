package dto

import "github.com/google/uuid"

type CreateProjectRequest struct {
	IdeaID      *uuid.UUID `json:"idea_id"`
	Title       string     `json:"title" binding:"required"`
	Description *string    `json:"description"`
}

type GetProjectsQuery struct {
	Status   *string    `form:"status"`
	LeaderID *uuid.UUID `form:"leader_id"`
	Search   *string    `form:"search"`
	StartAt  *int       `form:"start_at"`
	Limit    *int       `form:"limit"`
}
