package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
}
