package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

type User struct {
	InternalID int64          `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID      `json:"public_id" db:"public_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name       string         `json:"name" db:"name" gorm:"not null"`
	Email      string         `json:"email" db:"email" gorm:"unique;not null"`
	Password   string         `json:"-" db:"password" gorm:"not null"`
	Role       Role           `json:"role" db:"role" gorm:"type:varchar(50);not null;default:'user'"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
