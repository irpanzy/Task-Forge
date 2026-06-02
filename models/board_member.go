package models

import "time"

type BoardMember struct {
	BoardID  int64     `json:"board_id" db:"board_id" gorm:"primaryKey"`
	UserID   int64     `json:"user_id" db:"user_id" gorm:"primaryKey"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at" gorm:"autoCreateTime"`
}
