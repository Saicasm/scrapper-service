package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scraper/internal/controllers"
	"github.com/scraper/internal/handlers"
	"github.com/scraper/internal/middleware"
	"github.com/sirupsen/logrus"
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

func SetupControllerRoutes(router *gin.Engine, linkedInController *controllers.LinkedInController, log *logrus.Logger) {
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Replace with your allowed origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	v1 := router.Group("/api/v1")
	{
		jobs := v1.Group("/ingest")
		{
			jobs.POST("/:userId", linkedInController.CreateJob)
			jobs.GET("/:userId", linkedInController.GetJobsForUserID)
			jobs.GET("/analytics/:userId", linkedInController.GetAnalyticsForUser)
			jobs.PUT("/:jobId", linkedInController.UpdateJob)
			// Define other routes
		}
	}
}

// TODO: Create strcut and have a single Setup method
func SetupControllerRoutesForUser(router *gin.Engine, userController *controllers.UserController, log *logrus.Logger) {
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Replace with your allowed origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	v1 := router.Group("/api/v1")
	{
		jobs := v1.Group("/ingest/user")
		{
			jobs.POST("/", userController.Create)
			jobs.GET("/all", userController.GetAllUsers)
			jobs.PUT("/update/:userId", userController.UpdateUser)
			jobs.GET("/skills/:userId", userController.GetSkillsForUser)
			jobs.GET("/:userId", userController.GetUserById)
			// Define other routes
		}
	}
}
