package middleware

import (
	"life-gamifying/internal/utils"
	"net/http"

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
