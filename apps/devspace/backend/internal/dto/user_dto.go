package dto

import (
	"github.com/google/uuid"
)

type GetMeResponse struct {
	ID        uuid.UUID               `json:"id"`
	Email     string                  `json:"email"`
	Nickname  string                  `json:"nickname"`
	AvatarUri string                  `json:"avatar_uri"`
	Bio       string                  `json:"bio"`
	Skills    []SkillCategoryResponse `json:"skills"`
}

type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"min=3,max=50"`
	Bio      *string `json:"bio" binding:"min=3,max=255"`
}
