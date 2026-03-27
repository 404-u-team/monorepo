package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSlotRequest struct {
	PrimarySkillsID   []uuid.UUID `json:"primary_skills_id" binding:"required"`
	SecondarySkillsID []uuid.UUID `json:"secondary_skills_id" binding:"required"`
	Title             string      `json:"title" binding:"required,min=3,max=255"`
	Description       *string     `json:"description" binding:"omitempty,min=3,max=255"`
}

type UpdateSlotRequest struct {
	PrimarySkillsID   []uuid.UUID `json:"primary_skills_id" binding:"omitempty"`
	SecondarySkillsID []uuid.UUID `json:"secondary_skills_id" binding:"omitempty"`
	Title             *string     `json:"title" binding:"omitempty,min=3,max=255"`
	Description       *string     `json:"description" binding:"omitempty,min=3,max=255"`
	Status            *string     `json:"status" binding:"omitempty,oneof=open closed"`
}

type GetSlotResponse struct {
	ID              uuid.UUID               `json:"id"`
	PrimarySkills   []SkillCategoryResponse `json:"primary_skills"`
	SecondarySkills []SkillCategoryResponse `json:"secondary_skills"`
	Title           string                  `json:"title"`
	Description     *string                 `json:"description"`
	Status          string                  `json:"status"`
	UserID          *uuid.UUID              `json:"user_id"`
	CreatedAt       time.Time               `json:"created_at"`
}
