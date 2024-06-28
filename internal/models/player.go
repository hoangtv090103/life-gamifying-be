package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	ID           uint    `json:"id" gorm:"primaryKey"`
	UserID       uint    `json:"user_id" gorm:"unique;not null"`
	User         User    `json:"user" gorm:"foreignKey:UserID"`
	Level        uint    `json:"level"`
	RankID       uint    `json:"rank_id"`
	Rank         Rank    `json:"rank" gorm:"foreignKey:RankID"`
	Exp          uint    `json:"exp"`
	Health       uint    `json:"health"`
	Mana         uint    `json:"mana"`
	Strength     uint    `json:"strength"`
	Agility      uint    `json:"agility"`
	Intelligence uint    `json:"intelligence"`
	Habits       []Habit `json:"habits" gorm:"foreignKey:PlayerID"`
	Quests       []Quest `json:"quests" gorm:"foreignKey:PlayerID"`
	Dailies      []Daily `json:"dailies" gorm:"foreignKey:PlayerID"`
}
