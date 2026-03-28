package dto

import (
	"encoding/json"
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
	MainRole  uuid.UUID               `json:"main_role"`
	AvatarUri string                  `json:"avatar_uri"`
	Bio       string                  `json:"bio"`
	Skills    []SkillCategoryResponse `json:"skills"`
}

type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"min=3,max=50"`
	Bio      *string `json:"bio" binding:"min=3,max=255"`
}

type SafeUser struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string    `gorm:"unique; column:email; not null"`
	Nickname  string    `gorm:"column:nickname; not null"`
	AvatarUrl string    `gorm:"column:avatar_url; null"`
	MainRole  string    `gorm:"column:main_role; null"`
	Bio       string    `gorm:"column:bio; null"`
	IsAdmin   bool      `gorm:"column:is_admin; not null"`
}

type UUIDSlice []uuid.UUID

func (u *UUIDSlice) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 || string(bytes) == "null" {
		*u = nil
		return nil
	}

	var strings []string
	if err := json.Unmarshal(bytes, &strings); err != nil {
		return err
	}

	result := make(UUIDSlice, len(strings))
	for i, s := range strings {
		id, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		result[i] = id
	}

	*u = result
	return nil
}

func (u *UUIDSlice) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*u = nil
		return nil
	}

	// Если пришло как JSON массив
	if text[0] == '[' {
		var strings []string
		if err := json.Unmarshal(text, &strings); err != nil {
			return err
		}

		result := make(UUIDSlice, len(strings))
		for i, s := range strings {
			id, err := uuid.Parse(s)
			if err != nil {
				return err
			}
			result[i] = id
		}
		*u = result
		return nil
	}

	// Если пришло как один UUID
	id, err := uuid.Parse(string(text))
	if err != nil {
		return err
	}
	*u = UUIDSlice{id}
	return nil
}

func (u UUIDSlice) MarshalJSON() ([]byte, error) {
	if u == nil {
		return []byte("null"), nil
	}

	strings := make([]string, len(u))
	for i, id := range u {
		strings[i] = id.String()
	}
	return json.Marshal(strings)
}

type GetUsersRequest struct {
	StartAt  *uint      `form:"start_at" json:"start_at"`
	Limit    *uint      `form:"limit" json:"limit"`
	Username *string    `form:"username" json:"username"`
	MainRole *string    `form:"main_role" json:"main_role"`
	Skills   *UUIDSlice `form:"skills" json:"skills"`
}
