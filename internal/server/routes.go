package server

import (
	"life-gamifying/internal/middleware"
	"life-gamifying/internal/routes"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/cors"

)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// CORS
	// r.Use(cors.New(cors.Config{
	//     AllowOrigins:     []string{"https://foo.com"},
	//     AllowMethods:     []string{"PUT", "PATCH"},
	//     AllowHeaders:     []string{"Origin"},
	//     ExposeHeaders:    []string{"Content-Length"},
	//     AllowCredentials: true,
	//     AllowOriginFunc: func(origin string) bool {
	//       return origin == "https://github.com"
	//     },
	//     MaxAge: 12 * time.Hour,
	//   }))
	
	
	api := r.Group("/api")
	v1 := api.Group("/v1")
	routes.AuthRoutes(v1, s.db)
	
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
