package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint   `json:"id" gorm:"primaryKey"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Password  string `json:"password"`
}
