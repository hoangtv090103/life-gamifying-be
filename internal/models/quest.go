package models

import (
	"time"

	"gorm.io/gorm"
)

type Quest struct {
	gorm.Model
	ID            uint       `json:"id" gorm:"primaryKey"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Difficulty    uint       `json:"difficulty"`
	Completed     bool       `json:"completed"`
	PlayerID      uint       `json:"player_id"`
	Player        Player     `json:"player"`
	DueDate       *time.Time `json:"due_date"`
	CompletedDate *time.Time `json:"completed_date"`
}
