package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"-" gorm:"unique;not null"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
	Level  uint `json:"level" gorm:"default:1"`
	// RankID       uint    `json:"rank_id"`
	// Rank         Rank    `json:"rank" gorm:"foreignKey:RankID"`
	Exp          uint    `json:"exp" gorm:"default:0"`
	ExpToNext    uint    `json:"exp_to_next" gorm:"default:100"`
	Health       uint    `json:"health" gorm:"default:100"`
	Mana         uint    `json:"mana" gorm:"default:100"`
	Strength     uint    `json:"strength" gorm:"default:10"`
	Agility      uint    `json:"agility" gorm:"default:10"`
	Intelligence uint    `json:"intelligence" gorm:"default:10"`
	Habits       []Habit `json:"habits" gorm:"foreignKey:PlayerID"`
	Quests       []Quest `json:"quests" gorm:"foreignKey:PlayerID"`
	Dailies      []Daily `json:"dailies" gorm:"foreignKey:PlayerID"`
	Gold         uint    `json:"gold" gorm:"default:0"`
}
