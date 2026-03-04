package dto

type SkillCategoriesListRequest struct {
	ParentId *uint   `json:"parent_id"`
	Search   *string `json:"search"`
	Page     *uint   `json:"page"`
	Limit    *uint   `json:"limit"`
}
