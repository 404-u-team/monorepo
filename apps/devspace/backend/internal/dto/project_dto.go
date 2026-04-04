package dto

import (
	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	IdeaID      *uuid.UUID `json:"idea_id"`
	Title       string     `json:"title" binding:"required,min=3,max=255"`
	Description *string    `json:"description" binding:"omitempty,min=3,max=255"`
	Content     *string    `json:"content"`
}

type GetProjectsQuery struct {
	Status   *string    `form:"status" binding:"omitempty,oneof=open closed"`
	LeaderID *uuid.UUID `form:"leader_id"`
	Search   *string    `form:"search"`
	StartAt  *int       `form:"start_at"`
	Limit    *int       `form:"limit"`
}

type UpdateProjectRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=3,max=255"`
	Description *string `json:"description" binding:"omitempty,min=3,max=255"`
	Status      *string `json:"status" binding:"omitempty,oneof=open closed"`
	Content     *string `json:"content"`
}
