package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	ChatID    uuid.UUID `gorm:"not null" json:"chat_id"`
	SenderID  uuid.UUID `gorm:"not null" json:"sender_id"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
