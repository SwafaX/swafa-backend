package models

import (
	"github.com/google/uuid"
	"time"
)

type Item struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title       string    `gorm:"size:255;not null" json:"title" form:"title" binding:"required"`
	Description string    `gorm:"type:text" json:"description" form:"description"`
	Status      string    `gorm:"size:50;not null;default:'unfinished'" json:"status" form:"status"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ItemCreation struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
