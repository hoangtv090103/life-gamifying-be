package models

import "gorm.io/gorm"

type Rank struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"unique;not null"`
	MinLevel int    `json:"min_level" gorm:"not null"`
}
