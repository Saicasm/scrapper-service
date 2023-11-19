// cmd/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/scraper/constants"
	"github.com/scraper/internal/config"
	"github.com/scraper/internal/controllers"
	"github.com/scraper/internal/loggers"
	"github.com/scraper/internal/repositories"
	"github.com/scraper/internal/routes"
	"github.com/scraper/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	log := config.InitLogger()
	config.InitMongoDB()
	linkedInCollection := config.MongoDB.Collection("linkedin_jobs")

	// Create a LinkedIRepository as a value
	linkedInRepository := repositories.NewLinkedInRepository(linkedInCollection, log)
	// Create a LinkedInService and pass the repository
	linkedInService := services.NewLinkedInService(linkedInRepository, log)
	linkedInController := controllers.NewLinkedInController(linkedInService, log, "http://127.0.0.1:5000/api/v1/extractor/analyse")
	routes.SetupControllerRoutes(r, linkedInController, log)
	port := constants.ServerPort
	log.Printf("Server is running on port %s", port)
	logger.Info("Server is running on", logrus.Fields{"port": port})
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
