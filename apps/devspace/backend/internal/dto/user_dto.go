package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetMeResponse struct {
	ID        uuid.UUID               `json:"id"`
	Email     string                  `json:"email"`
	Nickname  string                  `json:"nickname"`
	AvatarUri string                  `json:"avatar_uri"`
	Bio       string                  `json:"bio"`
	CreatedAt time.Time               `json:"created_at"`
	Skills    []SkillCategoryResponse `json:"skills"`
}

type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"min=3,max=50"`
	Bio      *string `json:"bio" binding:"min=3,max=255"`
}
