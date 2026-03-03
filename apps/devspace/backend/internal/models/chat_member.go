package models

import (
	"time"
)

type ChatMember struct {
	ChatID     uint      `gorm:"column:chat_id; primaryKey"`
	UserID     uint      `gorm:"column:user_id; primaryKey"`
	JoinedAt   time.Time `gorm:"column:joined_at; not null"`
	LastReadAt time.Time `gorm:"column:last_read_at; not null"`

	Chat Chat `gorm:"foreignKey:ChatId"`
	User User `gorm:"foreignKey:UserId"`
}

func (cm *ChatMember) TableName() string { return "Chat_Member" }
