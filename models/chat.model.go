package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID           uuid.UUID  `gorm:"primaryKey" json:"id"`
	Participant1 uuid.UUID 	`gorm:"not null" json:"participant1_id"`
	Participant2 uuid.UUID  `gorm:"not null" json:"participant2_id"`
	CreatedAt    time.Time 	`json:"created_at"`
	UpdatedAt    time.Time 	`json:"updated_at"`
}
