package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Card struct {
	InternalID  int64          `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID    uuid.UUID      `json:"public_id" db:"public_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	ListID      int64          `json:"list_id" db:"list_id" gorm:"not null;index"`
	Title       string         `json:"title" db:"title" gorm:"not null"`
	Description string         `json:"description" db:"description"`
	Position    int            `json:"position" db:"position" gorm:"not null"`
	DueDate     *time.Time     `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
