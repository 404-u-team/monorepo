package dto

import "github.com/google/uuid"

type CreateSlotRequest struct {
	SkillCategoryID uuid.UUID `json:"skill_category_id" binding:"required"`
	Title           string    `json:"title" binding:"required,min=3,max=255"`
	Description     *string   `json:"description" binding:"omitempty,min=3,max=255"`
}

type UpdateSlotRequest struct {
	SkillCategoryID *uuid.UUID `json:"skill_category_id"`
	Title           *string    `json:"title" binding:"omitempty,min=3,max=255"`
	Description     *string    `json:"description" binding:"omitempty,min=3,max=255"`
	Status          *string    `json:"status" binding:"omitempty,oneof=open closed"`
}
