package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	InternarlID int64          `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID    uuid.UUID      `json:"public_id" db:"public_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string         `json:"name" db:"name" gorm:"not null"`
	Email       string         `json:"email" db:"email" gorm:"unique;not null"`
	Password    string         `json:"password" db:"password" gorm:"not null"`
	Role        string         `json:"role" db:"role" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:- gorm:"index"`
}
