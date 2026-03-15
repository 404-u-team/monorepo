package dto

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	IdeaID      *uuid.UUID `json:"idea_id"`
	Title       string     `json:"title" binding:"required,min=3,max=255"`
	Description *string    `json:"description" binding:"omitempty,min=3,max=255"`
}

type GetProjectsQuery struct {
	Status   *string    `form:"status" binding:"omitempty,oneof=open closed"`
	LeaderID *uuid.UUID `form:"leader_id"`
	Search   *string    `form:"search"`
	StartAt  *int       `form:"start_at"`
	Limit    *int       `form:"limit"`
}

type UpdateProjectRequest struct {
	Title       *string `json:"title" binding:"min=3,max=255"`
	Description *string `json:"description" binding:"min=3,max=255"`
	Status      *string `json:"status" binding:"omitempty,oneof=open closed"`
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

func FromUser(u *models.User) SafeUser {
	return SafeUser{
		ID:        u.ID,
		Email:     u.Email,
		Nickname:  u.Nickname,
		AvatarUrl: u.AvatarUrl,
		MainRole:  u.MainRole,
		Bio:       u.Bio,
		IsAdmin:   u.IsAdmin,
	}
}

type SafeRequest struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	SlotID      uuid.UUID `gorm:"column:slot_id;type:uuid; not null"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid; not null"`
	Type        string    `gorm:"column:type; not null"`
	Status      string    `gorm:"column:status; not null"`
	CoverLetter string    `gorm:"column:cover_letter; not null"`

	Slot models.ProjectSlot `gorm:"foreignKey:SlotID"`
	User SafeUser           `gorm:"foreignKey:UserID"`
}

func FromRequest(req *models.ProjectRequest) SafeRequest {
	return SafeRequest{
		ID:          req.ID,
		SlotID:      req.SlotID,
		UserID:      req.UserID,
		Type:        req.Type,
		Status:      req.Status,
		CoverLetter: req.CoverLetter,
		Slot:        req.Slot,
		User:        FromUser(&req.User),
	}
}

func FromRequests(requests []models.ProjectRequest) []SafeRequest {
	result := make([]SafeRequest, len(requests))
	for i, req := range requests {
		result[i] = FromRequest(&req)
	}
	return result
}
