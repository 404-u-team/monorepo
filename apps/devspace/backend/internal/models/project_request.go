package models

import "github.com/google/uuid"

type ProjectRequest struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	SlotID      uuid.UUID `gorm:"column:slot_id;type:uuid; not null"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid; not null"`
	Type        string    `gorm:"column:type; not null"`
	Status      string    `gorm:"column:status; not null"`
	CoverLetter string    `gorm:"column:cover_letter; not null"`

	Slot ProjectSlot `gorm:"foreignKey:SlotID"`
	User User        `gorm:"foreignKey:UserID"`
}

func (r *ProjectRequest) TableName() string { return "Project_Request" }
