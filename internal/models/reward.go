package models

import "gorm.io/gorm"

type Reward struct {
	gorm.Model
	ID          uint         `json:"id" gorm:"primaryKey"`
	PlayerID    uint         `json:"player_id"`
	Player      Player       `json:"player" gorm:"foreignKey:PlayerID"`
	RewardItems []RewardItem `json:"reward_items" gorm:"foreignKey:RewardID"`
}

type RewardItem struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	RewardID uint   `json:"reward_id"`
	Reward   Reward `json:"reward" gorm:"foreignKey:reward_id"`
}
