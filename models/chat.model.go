package models

import (
	"time"
)

type Chat struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Participant1 uint      `gorm:"not null" json:"participant1_id"`
	Participant2 uint      `gorm:"not null" json:"participant2_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
