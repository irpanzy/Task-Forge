package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type List struct {
	InternalID    int64          `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID      uuid.UUID      `json:"public_id" db:"public_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	BoardID       int64          `json:"board_id" db:"board_id" gorm:"not null;index"`
	BoardPublicID uuid.UUID      `json:"board_public_id" db:"board_public_id" gorm:"type:uuid;not null;index"`
	Title         string         `json:"title" db:"title" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
