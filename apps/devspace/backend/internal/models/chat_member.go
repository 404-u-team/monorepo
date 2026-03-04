package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatMember struct {
	ChatID     uuid.UUID `gorm:"column:chat_id;type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey"`
	JoinedAt   time.Time `gorm:"column:joined_at; not null"`
	LastReadAt time.Time `gorm:"column:last_read_at; not null"`

	Chat Chat `gorm:"foreignKey:ChatID"`
	User User `gorm:"foreignKey:UserID"`
}

func (cm *ChatMember) TableName() string { return "Chat_Member" }
