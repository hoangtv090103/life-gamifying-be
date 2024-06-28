package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int   `json:"price"`
	Armor       int   `json:"armor"`
	Strength    int   `json:"strength"`
	Agility     int   `json:"agility"`
	Intelligence int   `json:"intelligence"`
}
