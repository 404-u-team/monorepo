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

// Удалить/привязать навык одинаковы по содержанию, смысл дублировать?
type BaseSkillRequest struct {
	SkillID uuid.UUID `json:"skill_id" form:"skill_id"`
}
