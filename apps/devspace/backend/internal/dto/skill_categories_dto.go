package dto

import "github.com/google/uuid"

type SkillCategoriesListQuery struct {
	ParentId *uuid.UUID `form:"parent_id"`
	Search   *string    `form:"search"`
	StartAt  *int       `form:"start_at"`
	Limit    *int       `form:"limit"`
}

type SkillCategoryAddRequest struct {
	ParentId *uuid.UUID `json:"parent_id"`
	Name     string     `json:"name"`
}

type SkillCategoryResponse struct {
	ID       uuid.UUID               `json:"id"`
	ParentID *uuid.UUID              `json:"parent_id"`
	Name     string                  `json:"name"`
	Icon     string                  `json:"icon"`
	Color    string                  `json:"color"`
	Children []SkillCategoryResponse `json:"children"`
}

// Удалить/привязать навык одинаковы по содержанию, смысл дублировать?
type BaseSkillRequest struct {
	SkillID uuid.UUID `json:"skill_id" form:"skill_id"`
}
