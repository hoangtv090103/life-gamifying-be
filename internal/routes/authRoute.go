package routes

import (
	"life-gamifying/internal/database"
	"life-gamifying/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(group *gin.RouterGroup, s database.Service) {
	authRoute := group.Group("/auth")
	{
		authRoute.POST("/login", func(c *gin.Context) {
			handlers.Login(c, s)
		})
		authRoute.POST("/register", func(c *gin.Context) {
			handlers.Register(c, s)
		})
	}
}
