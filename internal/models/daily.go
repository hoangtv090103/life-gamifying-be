package models

import (
	"time"
	"gorm.io/gorm"
)

type Daily struct {
	gorm.Model
	ID            uint       `json:"id" gorm:"primaryKey"`
	PlayerID      uint       `json:"player_id"`
	Player        Player     `json:"player" gorm:"foreignKey:PlayerID"`
	Completed     bool       `json:"completed"`
	RankID        uint       `json:"-"`
	Rank          Rank       `json:"rank" gorm:"foreignKey:RankID"`
	DueDate       *time.Time `json:"due_date"`
	CompletedDate *time.Time `json:"completed_date"`
}
