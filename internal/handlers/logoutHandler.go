package handlers

import (
	"life-gamifying/internal/database"
	"life-gamifying/internal/models"
	"life-gamifying/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(ctx *gin.Context, s database.Service) {
	// Get the token from the header
	token := ctx.GetHeader("Authorization")

	// Check if token is empty
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		ctx.Abort()
		return
	}

	// Check if token is valid
	_, err := utils.ValidateToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	// Delete token from database
	db := s.DB()
	db.Model(&models.Token{}).Where("token = ?", token).Delete(&models.Token{})

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successfully"})
}
