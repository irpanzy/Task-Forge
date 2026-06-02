package models

type CardAssignee struct {
	CardID int64 `json:"card_id" db:"card_id" gorm:"not null;index;uniqueIndex:idx_card_assignee"`
	UserID int64 `json:"user_id" db:"user_id" gorm:"not null;index;uniqueIndex:idx_card_assignee"`
}
