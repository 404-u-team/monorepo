package models

import (
	"time"
)

type Message struct {
	Id       uint      `gorm:"primaryKey; column:id"`
	ChatId   uint      `gorm:"column:chat_id; not null"`
	SenderId uint      `gorm:"column:sender_id; not null"`
	Content  string    `gorm:"column:content; not null"`
	SentAt   time.Time `gorm:"column:sent_at; not null"`
	IsEdited bool      `gorm:"column:is_edited; not null"`

	Chat   Chat `gorm:"foreignKey:ChatId"`
	Sender User `gorm:"foreignKey:SenderId"`
}

func (m *Message) TableName() string { return "Message" }
