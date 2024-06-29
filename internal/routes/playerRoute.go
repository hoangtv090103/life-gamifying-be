package routes

import (
	"life-gamifying/internal/database"
	"life-gamifying/internal/handlers"

	"github.com/gin-gonic/gin"
)

func PlayerRoutes(group *gin.RouterGroup, s database.Service) {
	playerRoute := group.Group("/players")
	{
		playerRoute.GET("/", func(c *gin.Context) {
			handlers.GetPlayers(c, s)
		})
		playerRoute.GET("/:id", func(c *gin.Context) {
			handlers.GetPlayerByID(c, s)
		})
		playerRoute.POST("/", func(c *gin.Context) {
			handlers.CreatePlayer(c, s)
		})

		playerRoute.PUT("/:id", func(c *gin.Context) {
			handlers.UpdatePlayer(c, s)
		})

		playerRoute.DELETE("/:id", func(c *gin.Context) {
			handlers.DeletePlayer(c, s)
		})

		playerRoute.GET("/:id/habits", func(c *gin.Context) {
			handlers.GetHabitsByPlayerID(c, s)
		})

	}
}
