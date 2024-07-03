package handlers

import (
	"life-gamifying/internal/database"
	"life-gamifying/internal/models"
	"time"
)

func SaveToken(s database.Service, token string, uid uint) error {
	db := s.DB()
	expireAt := time.Now().Add(24 * time.Hour)
	tokenObj := models.Token{
		UserID:   uid,
		Token:    token,
		ExpireAt: expireAt,
	}

	err := db.Create(&tokenObj).Error

	if err != nil {
		return err
	}

	return nil
}

