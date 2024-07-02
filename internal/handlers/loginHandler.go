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

	// Check if username or email is empty
	var login string

	if loginUser.Username != "" {
		login = loginUser.Username
	} else if loginUser.Email != "" {
		login = loginUser.Email
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username/Email is required"})
		return err
	}

	password := loginUser.Password

	// Check if username or email is empty
	if password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return err
	}

	db := s.DB()
	// Check if user exists by email or username
	var user models.User
	db.Where("email = ? OR username = ?", login, login).First(&user)

	// Check if password is correct
	if !utils.CheckPasswordHash(password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return err
	}

	token, err := utils.CreateToken(user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	var player models.Player
	db.Where("user_id = ?", user.ID).First(&player)

	ctx.JSON(http.StatusOK, gin.H{"token": token, "player_id": player.ID, "message": "Successfully logged in"})

	return nil
}
