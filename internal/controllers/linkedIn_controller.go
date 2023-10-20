// server/controllers/todo_controller.go
package controllers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/scraper/internal/models"
	"github.com/scraper/internal/services"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

type LinkedInController struct {
	Service services.LinkedInService
	Log     *logrus.Logger
}

func NewLinkedInController(service services.LinkedInService, log *logrus.Logger) *LinkedInController {
	return &LinkedInController{
		Service: service,
		Log:     log,
	}
}

func (c *LinkedInController) CreateJob(ctx *gin.Context) {
	// Read the request body
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	htmlContent := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse HTML content"})
		return
	}

	linkedin := models.NewLinkedIn()

	// Create channels for extracted data
	titleChan := make(chan string)
	descriptionChan := make(chan string)
	skillsChan := make(chan string)

	// Launch Goroutines to extract data from different elements
	go extractData(doc, ".job-details-jobs-unified-top-card__job-title", titleChan)
	go extractData(doc, ".job-details-jobs-unified-top-card__primary-description", descriptionChan)
	go extractData(doc, ".job-details-how-you-match__skills-item-subtitle", skillsChan)

	// Receive extracted data from channels
	linkedin.Title = <-titleChan
	linkedin.JobDescription = <-descriptionChan
	linkedin.Skills = <-skillsChan

	// Close the channels
	//close(titleChan)
	//close(descriptionChan)
	//close(skillsChan)

	// Clean the data
	cleanData(&linkedin.Title)
	cleanData(&linkedin.JobDescription)
	cleanData(&linkedin.Skills)

	if err := c.Service.Create(&linkedin); err != nil {
		c.Log.WithError(err).Error("Failed to create a new job")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new job"})
	} else {
		ctx.JSON(http.StatusCreated, linkedin)
	}
}

func extractData(doc *goquery.Document, selector string, dataChan chan string) {
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		dataChan <- s.Text()

	})
}

func cleanData(field *string) {
	*field = strings.TrimSpace(*field)            // Remove leading and trailing white spaces
	*field = strings.ReplaceAll(*field, "\n", "") // Remove newline characters
}

// Implement other CRUD operations in a similar manner.
