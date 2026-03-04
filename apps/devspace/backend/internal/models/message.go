package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ChatId   uuid.UUID `gorm:"column:chat_id;type:uuid; not null"`
	SenderId uuid.UUID `gorm:"column:sender_id;type:uuid; not null"`
	Content  string    `gorm:"column:content; not null"`
	SentAt   time.Time `gorm:"column:sent_at; not null"`
	IsEdited bool      `gorm:"column:is_edited; not null"`

	Chat   Chat `gorm:"foreignKey:ChatId"`
	Sender User `gorm:"foreignKey:SenderId"`
}

func (m *Message) TableName() string { return "Message" }
