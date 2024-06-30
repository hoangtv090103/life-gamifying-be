package server

import (
	"life-gamifying/internal/middleware"
	"life-gamifying/internal/routes"
	"life-gamifying/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)
func (s *Server) RegisterRoutes() http.Handler {
    r := gin.Default()

    // Apply CORS middleware globally
    r.Use(utils.CORS())

    // Define routes
    api := r.Group("/api")
    api.GET("/health/redis", s.redisHealthHandler)
    api.GET("/health/postgres", s.postgresHealthHandler)
    routes.AuthRoutes(api, s.db)

    v1 := api.Group("/v1")
    v1.Use(middleware.AuthMiddleware())
    routes.HabitRoutes(v1, s.db)
    routes.UserRoutes(v1, s.db)
    routes.PlayerRoutes(v1, s.db)

    r.GET("/", s.HelloWorldHandler)

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
