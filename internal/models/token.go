package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	*gorm.Model
	UserID   uint      `json:"user_id"`
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}
