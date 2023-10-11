// cmd/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/scraper/constants"
	"github.com/scraper/internal/routes"
	"log"
)

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)

	// Start the server
	port := constants.ServerPort
	log.Printf("Server is running on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
