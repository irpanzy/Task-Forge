package models

import (
	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/models/types"
)

type ListPosition struct {
	InternalID int64           `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID       `json:"public_id" db:"public_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	BoardID    int64           `json:"board_id" db:"board_id" gorm:"not null;index"`
	ListOrder  types.UUIDArray `json:"list_order" db:"list_order" gorm:"type:uuid[];not null"`
}
