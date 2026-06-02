package models

type CardLabel struct {
	CardID  int64 `json:"card_id" db:"card_id" gorm:"not null;index"`
	LabelID int64 `json:"label_id" db:"label_id" gorm:"not null;index"`
}
