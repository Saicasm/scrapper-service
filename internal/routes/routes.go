package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scraper/internal/handlers"
	"github.com/scraper/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Replace with your allowed origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))
	// API routes
	api := r.Group("/api")
	{
		tasks := api.Group("/ingest")
		{
			tasks.GET("/health", handlers.Health)
			tasks.POST("/", handlers.IngestData)
		}
	}
}
