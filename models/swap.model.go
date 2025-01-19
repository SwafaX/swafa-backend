package models

import (
	"github.com/google/uuid"
	"time"
)

type Swap struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	RequesterID   uuid.UUID `gorm:"type:uuid;not null" json:"requester_id"`    // User who initiated the swap
	RecipientID   uuid.UUID `gorm:"type:uuid;not null" json:"recipient_id"`    // User who owns the requested item
	RequestItemID uuid.UUID `gorm:"type:uuid;not null" json:"request_item_id"` // Item the requester wants
	// 'pending', 'accepted', 'rejected', 'cancelled', 'completed'
	Status    string    `gorm:"size:50;not null;default:'pending'" json:"status"`
	Message   string    `gorm:"type:text" json:"message,omitempty"` // Optional message from the requester
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
