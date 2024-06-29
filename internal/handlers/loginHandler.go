package handlers

import (
	"life-gamifying/internal/database"
	"life-gamifying/internal/models"
	"life-gamifying/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context, s database.Service) error {
	var loginUser models.User
	err := ctx.ShouldBindJSON(&loginUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	db := s.DB()
	// Check if user exists by email or username
	var user models.User
	db.Where("email = ? OR username = ?", user.Email, user.Username).First(&user)

	// Check if password is correct
	if !utils.CheckPasswordHash(loginUser.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return err
	}

	ctx.JSON(http.StatusOK, user)

	return nil
}
