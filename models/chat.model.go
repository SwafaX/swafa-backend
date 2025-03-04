package models

import (
	"github.com/google/uuid"
	"time"
)

type Chat struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	SwapID    uuid.UUID `gorm:"type:uuid;not null" json:"swap_id,omitempty"`
	User1     uuid.UUID `gorm:"type:uuid;not null" json:"user1_id,omitempty"`
	User2     uuid.UUID `gorm:"type:uuid;not null" json:"user2_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
