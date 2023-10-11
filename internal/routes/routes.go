package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/scraper/internal/handlers"
	"github.com/scraper/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())

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
