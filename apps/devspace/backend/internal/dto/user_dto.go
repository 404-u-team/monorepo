package dto

import (
	"time"

	"github.com/google/uuid"
)

type PrivateUserProfile struct {
	ID        uuid.UUID               `json:"id"`
	Email     string                  `json:"email"`
	Nickname  string                  `json:"nickname"`
	MainRole  string                  `json:"main_role"`
	AvatarUri string                  `json:"avatar_uri"`
	Bio       string                  `json:"bio"`
	Skills    []SkillCategoryResponse `json:"skills"`
	CreatedAt time.Time               `json:"created_at"`
}

type PublicUserProfile struct {
	ID        uuid.UUID               `json:"id"`
	Nickname  string                  `json:"nickname"`
	MainRole  string                  `json:"main_role"`
	AvatarUri string                  `json:"avatar_uri"`
	Bio       string                  `json:"bio"`
	Skills    []SkillCategoryResponse `json:"skills"`
}

type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"min=3,max=50"`
	Bio      *string `json:"bio" binding:"min=3,max=255"`
}
