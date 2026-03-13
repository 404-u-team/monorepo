package dto

import "github.com/google/uuid"

type SkillCategoriesListRequest struct {
	ParentId *uuid.UUID `form:"parent_id"`
	Search   *string    `form:"search"`
	Page     *uint      `form:"page"`
	Limit    *uint      `form:"limit"`
}

type SkillCategoryAddRequest struct {
	ParentId *string `json:"parent_id"`
	Name     string  `json:"name"`
}

type SkillCategoryResponse struct {
	ID       uuid.UUID               `json:"id"`
	Name     string                  `json:"name"`
	ParentID *uuid.UUID              `json:"parent_id"`
	Children []SkillCategoryResponse `json:"children"`
}
