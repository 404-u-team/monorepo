package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectRequest struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SlotID      uuid.UUID `gorm:"column:slot_id;type:uuid; not null;uniqueIndex:idx_slot_user" json:"slot_id"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid; not null;uniqueIndex:idx_slot_user" json:"user_id"`
	Type        string    `gorm:"column:type; not null" json:"type"`
	Status      string    `gorm:"column:status; not null" json:"status"`
	CoverLetter string    `gorm:"column:cover_letter; not null" json:"cover_letter"`
	CreatedAt   time.Time `gorm:"column:created_at; not null" json:"created_at"`

	Slot ProjectSlot `gorm:"foreignKey:SlotID"`
	User User        `gorm:"foreignKey:UserID"`
}

func (r *ProjectRequest) TableName() string { return "Project_Request" }
