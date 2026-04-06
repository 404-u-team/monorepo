package dto

import (
	"encoding/json"
	"time"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
)

type PrivateUserProfile struct {
	ID        uuid.UUID               `json:"id"`
	Email     string                  `json:"email"`
	Nickname  string                  `json:"nickname"`
	MainRole  *models.SkillCategory   `json:"main_role"`
	AvatarUrl string                  `json:"avatar_url"`
	Bio       string                  `json:"bio"`
	Skills    []SkillCategoryResponse `json:"skills"`
	CreatedAt time.Time               `json:"created_at"`
}

type PublicUserProfile struct {
	ID        uuid.UUID               `json:"id"`
	Nickname  string                  `json:"nickname"`
	MainRole  *models.SkillCategory   `json:"main_role"`
	AvatarUrl string                  `json:"avatar_url"`
	Bio       string                  `json:"bio"`
	Skills    []SkillCategoryResponse `json:"skills"`
}

type GetUsersResponse struct {
	Total    int64               `json:"total"`
	Profiles []PublicUserProfile `json:"profiles"`
}

type UpdateUserRequest struct {
	Nickname  *string      `json:"nickname" binding:"omitempty,min=3,max=50"`
	MainRole  OptionalUUID `json:"main_role"`
	AvatarUrl *string      `json:"avatar_url"`
	Bio       *string      `json:"bio" binding:"omitempty,min=3,max=255"`
}

// OptionalUUID distinguishes between omitted field and explicit null in JSON.
type OptionalUUID struct {
	IsSet bool
	Value *uuid.UUID
}

func (o *OptionalUUID) UnmarshalJSON(data []byte) error {
	o.IsSet = true

	if string(data) == "null" {
		o.Value = nil
		return nil
	}

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	parsed, err := uuid.Parse(raw)
	if err != nil {
		return err
	}

	o.Value = &parsed
	return nil
}

type GetUsersRequest struct {
	StartAt  *uint      `form:"start_at" json:"start_at"`
	Limit    *uint      `form:"limit" json:"limit"`
	Search   *string    `form:"search" json:"search"`
	MainRole *string    `form:"main_role" json:"main_role"`
	Skills   *UUIDSlice `form:"skills" json:"skills"`
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

func (u *UUIDSlice) UnmarshalParam(param string) error {
	return u.UnmarshalText([]byte(param))
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
