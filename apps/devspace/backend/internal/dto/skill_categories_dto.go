package dto

type SkillCategoriesListRequest struct {
	ParentId *string `json:"parent_id"`
	Search   *string `json:"search"`
	Page     *uint   `json:"page"`
	Limit    *uint   `json:"limit"`
}

type SkillCategoryAddRequest struct {
	ParentId *string `json:"parent_id"`
	Name     string  `json:"name"`
}
