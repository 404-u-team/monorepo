package dto

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
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

type GetProjectsResponse struct {
	Total    int64            `json:"total"`
	Projects []models.Project `json:"projects"`
}

// QueryUUID supports form/query binding when value is either plain UUID
// or JSON array encoded by binder (for example: ["uuid"]).
type QueryUUID uuid.UUID

func (u *QueryUUID) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	raw := strings.TrimSpace(string(text))

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
