package models

import (
	"time"
)

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    uint      `gorm:"not null" json:"chat_id"`
	SenderID  uint      `gorm:"not null" json:"sender_id"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
