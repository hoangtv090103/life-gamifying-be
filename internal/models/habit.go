package models

import (
	"life-gamifying/internal/utils"

	"gorm.io/gorm"
)

type Habit struct {
	gorm.Model
	ID          uint             `json:"id" gorm:"primaryKey"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Difficulty  utils.Difficulty `json:"difficulty"`
	Frequency   utils.Frequency  `json:"frequency"`
	Success     uint             `json:"success"`
	Failure     uint             `json:"failure"`
	PlayerID    uint             `json:"player_id"`
	Player      Player           `json:"-" gorm:"foreignKey:PlayerID"`
}
