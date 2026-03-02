package models

import (
	"time"
)

type Chat_Member struct {
	ChatId     uint      `gorm:"column:chat_id; primaryKey"`
	UserId     uint      `gorm:"column:user_id; primaryKey"`
	JoinedAt   time.Time `gorm:"column:joined_at; not null"`
	LastReadAt time.Time `gorm:"column:last_read_at; not null"`

	Chat Chat `gorm:"foreignKey:ChatId"`
	User User `gorm:"foreignKey:UserId"`
}

func (cm *Chat_Member) TableName() string { return "Chat_Member" }
