package server

import (
	"life-gamifying/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		v1.GET("/habits", func(c *gin.Context) {
			handlers.GetAllHabits(c, s.db)
		})
		v1.GET("/habits/:id", func(c *gin.Context) {
			handlers.GetHabitByID(c, s.db)
		})
		v1.GET("/players/:id/habits", func(c *gin.Context) {
			handlers.GetHabitsOfPlayer(c, s.db)
		})
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) redisHealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.RDHealth())
}

func (s *Server) postgresHealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.PHealth())
}
