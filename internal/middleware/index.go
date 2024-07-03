package middleware

import (
	"life-gamifying/internal/database"
	"life-gamifying/internal/models"
	"life-gamifying/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		ctx.Next()
	}
}

func AuthenticateRequest(s database.Service) gin.HandlerFunc {
	db := s.DB()
	return func(ctx *gin.Context) {
		// Lấy token từ header hoặc cookie
		tokenStr := ctx.GetHeader("Authorization")
		var token models.Token
		err := db.Model(&models.Token{}).Where("token = ?", tokenStr).First(&token).Error

		// Kiểm tra lỗi
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Kiểm tra token hết hạn
		if token.ExpireAt.Before(time.Now()) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
