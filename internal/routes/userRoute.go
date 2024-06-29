package routes

import (
	"life-gamifying/internal/database"
	"life-gamifying/internal/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(group *gin.RouterGroup, s database.Service) {
	userRoute := group.Group("/users")
	{
		userRoute.GET("/", func(c *gin.Context) {
			handlers.GetUsers(c, s)
		})
		userRoute.GET("/:id", func(c *gin.Context) {
			handlers.GetUserByID(c, s)
		})
		userRoute.POST("/", func(c *gin.Context) {
			handlers.CreateUser(c, s)
		})

		userRoute.PUT("/:id", func(c *gin.Context) {
			handlers.UpdateUser(c, s)
		})

		userRoute.DELETE("/:id", func(c *gin.Context) {
			handlers.DeleteUser(c, s)
		})

	}
}
