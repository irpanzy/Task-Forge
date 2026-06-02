package models

import (
	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/models/types"
)

type CardPosition struct {
	InternalID int64           `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID       `json:"public_id" db:"public_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	ListID     int64           `json:"list_id" db:"list_id" gorm:"not null;index"`
	CardOrder  types.UUIDArray `json:"card_order" db:"card_order" gorm:"type:uuid[];not null"`
}
