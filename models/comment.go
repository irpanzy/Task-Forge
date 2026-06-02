package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	InternalID   int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID     uuid.UUID `json:"public_id" db:"public_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CardID       int64     `json:"card_id" db:"card_id" gorm:"not null;index"`
	CardPublicID uuid.UUID `json:"card_public_id" db:"card_public_id" gorm:"type:uuid;not null;index"`
	UserID       int64     `json:"user_id" db:"user_id" gorm:"not null;index"`
	UserPublicID uuid.UUID `json:"user_public_id" db:"user_public_id" gorm:"type:uuid;not null;index"`
	Message      string    `json:"message" db:"message" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}
