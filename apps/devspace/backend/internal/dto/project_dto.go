package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	IdeaID      *uuid.UUID `json:"idea_id"`
	Title       string     `json:"title" binding:"required,min=3,max=255"`
	Description *string    `json:"description" binding:"omitempty,min=3,max=255"`
	Content     *string    `json:"content"`
}

type GetProjectsQuery struct {
	Status      *string    `form:"status" binding:"omitempty,oneof=open closed"`
	LeaderID    *QueryUUID `form:"leader_id" binding:"omitempty"`
	Search      *string    `form:"search"`
	IdeaID      *QueryUUID `form:"idea_id"`
	OpenSlots   *bool      `form:"open_slots"`
	SlotsSkills *UUIDSlice `form:"slots_skills"`
	MinPeople   *int       `form:"min_people" binding:"omitempty,min=1,max=255"`
	MaxPeople   *int       `form:"max_people" binding:"omitempty,min=1,max=255"`
	StartAt     *int       `form:"start_at"`
	Limit       *int       `form:"limit"`
}

type UpdateProjectRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=3,max=255"`
	Description *string `json:"description" binding:"omitempty,min=3,max=255"`
	Status      *string `json:"status" binding:"omitempty,oneof=open closed"`
	Content     *string `json:"content"`
}

type ProjectBlock struct {
	ID             uuid.UUID  `json:"id"`
	LeaderID       uuid.UUID  `json:"leader_id"`
	IsLeader       bool       `json:"is_leader"`
	IsFavorite     bool       `json:"is_favorite"`
	Title          string     `json:"title"`
	Description    *string    `json:"description"`
	ViewsCount     uint       `json:"views_count"`
	FavoritesCount uint       `json:"favorites_count"`
	Status         string     `json:"status"`
	IdeaID         *uuid.UUID `json:"idea_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type GetProjectsResponse struct {
	Total    int64          `json:"total"`
	Projects []ProjectBlock `json:"projects"`
}

type GetProjectResponse struct {
	ID             uuid.UUID         `gorm:"column:id" json:"id"`
	LeaderID       uuid.UUID         `gorm:"column:leader_id" json:"leader_id"`
	IsLeader       bool              `gorm:"column:is_leader" json:"is_leader"`
	IsFavorite     bool              `gorm:"column:is_favorite" json:"is_favorite"`
	Title          string            `gorm:"column:title" json:"title"`
	Description    *string           `gorm:"column:description" json:"description"`
	Content        *string           `gorm:"column:content" json:"content"`
	ViewsCount     uint              `gorm:"column:views_count" json:"views_count"`
	FavoritesCount uint              `gorm:"column:favorites_count" json:"favorites_count"`
	Status         string            `gorm:"column:status" json:"status"`
	IdeaID         *uuid.UUID        `gorm:"column:idea_id" json:"idea_id"`
	CreatedAt      time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time         `gorm:"column:updated_at" json:"updated_at"`
	Slots          []GetSlotResponse `gorm:"-" json:"slots"`
}

// QueryUUID supports form/query binding when value is either plain UUID
// or JSON array encoded by binder (for example: ["uuid"]).
type QueryUUID uuid.UUID

func (u *QueryUUID) UnmarshalText(text []byte) error {
	raw := strings.TrimSpace(string(text))
	if raw == "" {
		return fmt.Errorf("uuid cannot be empty")
	}

	// 1) Plain UUID
	if parsed, err := uuid.Parse(raw); err == nil {
		*u = QueryUUID(parsed)
		return nil
	}

	// 2) JSON array: ["uuid"]
	if strings.HasPrefix(raw, "[") {
		var values []string
		if err := json.Unmarshal([]byte(raw), &values); err != nil {
			return err
		}
		if len(values) == 0 {
			return fmt.Errorf("empty uuid array")
		}
		raw = values[0]
		parsed, err := uuid.Parse(raw)
		if err != nil {
			return err
		}
		*u = QueryUUID(parsed)
		return nil
	}

	// 3) JSON string containing array: "[\"uuid\"]"
	if strings.HasPrefix(raw, "\"") {
		var encoded string
		if err := json.Unmarshal([]byte(raw), &encoded); err == nil {
			encoded = strings.TrimSpace(encoded)
			if strings.HasPrefix(encoded, "[") {
				var values []string
				if err := json.Unmarshal([]byte(encoded), &values); err != nil {
					return err
				}
				if len(values) == 0 {
					return fmt.Errorf("empty uuid array")
				}
				parsed, err := uuid.Parse(values[0])
				if err != nil {
					return err
				}
				*u = QueryUUID(parsed)
				return nil
			}

			parsed, err := uuid.Parse(encoded)
			if err != nil {
				return err
			}
			*u = QueryUUID(parsed)
			return nil
		}
	}

	parsed, err := uuid.Parse(raw)
	if err != nil {
		return err
	}

	*u = QueryUUID(parsed)
	return nil
}

func (u *QueryUUID) UnmarshalParam(param string) error {
	return u.UnmarshalText([]byte(param))
}

func (u QueryUUID) UUID() uuid.UUID {
	return uuid.UUID(u)
}
