package models

import (
	"time"

	"github.com/google/uuid"
)

type CardAttachment struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CardID     int64     `json:"card_id" db:"card_id" gorm:"not null;index"`
	UserID     int64     `json:"user_id" db:"user_id" gorm:"not null;index"`
	File       string    `json:"file" db:"file" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}
