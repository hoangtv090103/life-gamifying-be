package server

import (
	"life-gamifying/internal/middleware"
	"life-gamifying/internal/routes"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowWildcard = true
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type", "Content-Length"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	api := r.Group("/api")

	api.GET("/health/redis", s.redisHealthHandler)
	api.GET("/health/postgres", s.postgresHealthHandler)

	v1 := api.Group("/v1")
	routes.AuthRoutes(v1, s.db)

	v1.Use(cors.New(config))
	v1.Use(middleware.AuthMiddleware())

	r.GET("/", s.HelloWorldHandler)

	routes.HabitRoutes(v1, s.db)
	routes.UserRoutes(v1, s.db)
	routes.PlayerRoutes(v1, s.db)

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
