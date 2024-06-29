package routes

import (
	"life-gamifying/internal/database"
	"life-gamifying/internal/handlers"

	"github.com/gin-gonic/gin"
)

func HabitRoutes(group *gin.RouterGroup, s database.Service) {	
	habitRoute := group.Group("/habits")
	{
		habitRoute.GET("/", func(c *gin.Context) {
			handlers.GetAllHabits(c, s)
		})
		habitRoute.GET("/:id", func(c *gin.Context) {
			handlers.GetHabitByID(c, s)
		})
		habitRoute.POST("/", func(c *gin.Context) {
			handlers.CreateHabit(c, s)
		})

		habitRoute.PUT("/:id", func(c *gin.Context) {
			handlers.UpdateHabit(c, s)
		})

		habitRoute.DELETE("/:id", func(c *gin.Context) {
			handlers.DeleteHabit(c, s)
		})

	}
}
