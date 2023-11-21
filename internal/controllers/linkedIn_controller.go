// server/controllers/todo_controller.go
package controllers

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/scraper/internal/models"
	"github.com/scraper/internal/services"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type LinkedInController struct {
	Service       services.LinkedInService
	Log           *logrus.Logger
	SkillEndpoint string // Define the endpoint URL for skills
}

func NewLinkedInController(service services.LinkedInService, log *logrus.Logger, skillEndpoint string) *LinkedInController {
	return &LinkedInController{
		Service:       service,
		Log:           log,
		SkillEndpoint: skillEndpoint, // Initialize the endpoint URL
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
	compensationChan := make(chan string)
	companyChan := make(chan string)

	// Launch Goroutines to extract data from different elements
	// TODO: Extract skills from the description  chan
	go extractData(doc, ".job-details-jobs-unified-top-card__job-title", titleChan)
	go extractData(doc, ".jobs-description__content", descriptionChan)
	go extractData(doc, ".jobs-description__content", skillsChan)
	go extractData(doc, ".job-details-jobs-unified-top-card__job-insight", compensationChan)
	go extractData(doc, ".job-details-jobs-unified-top-card__primary-description", companyChan)
	// Receive extracted data from channels

	// Store the data into the structs
	linkedin.Title = <-titleChan
	linkedin.JobDescription = <-descriptionChan
	linkedin.Skills = <-skillsChan
	linkedin.Compensation = <-compensationChan
	linkedin.CompanyName = <-companyChan

	// Clean the data
	cleanData(&linkedin.Title)
	cleanData(&linkedin.JobDescription)
	cleanData(&linkedin.Skills)
	cleanData(&linkedin.Compensation)
	fmt.Println("Compensation", linkedin.Compensation)

	extractedSkills := linkedin.Skills
	cleanedText := removeWhitespaces(linkedin.Compensation)

	cleanedComp := extractCompensation(cleanedText)
	fmt.Println("Cleamned one: ", cleanedComp)

	linkedin.Compensation = cleanedComp
	newSkills, err := c.sendSkillsToEndpoint(extractedSkills)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send skills data"})
		return
	}

	// Update the LinkedIn model with the new skills
	linkedin.Skills = newSkills
	splitNameFromJD(&linkedin)
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
func extractCompensation(cleanedText string) string {
	// Assuming the compensation is in the format "$X,XXX/month"

	compensationHyphenIndex := strings.Index(cleanedText, "-")
	if compensationHyphenIndex != -1 && (cleanedText[(compensationHyphenIndex-1)] == 'r') {
		fmt.Println("Cleaned with hyphen", cleanedText[:(compensationHyphenIndex+(compensationHyphenIndex-1)+3)])
		return cleanedText[:(compensationHyphenIndex + (compensationHyphenIndex - 1) + 3)]
	} else {
		compensationRIndex := strings.IndexByte(cleanedText, 'r')

		fmt.Println("Cleaned COmp", cleanedText[:compensationRIndex+1])
		return cleanedText[:compensationRIndex+1]
	}
}
func cleanData(field *string) {
	*field = strings.ReplaceAll(*field, "\n", " ") // Remove newline characters
	*field = strings.TrimSpace(*field)             // Remove leading and trailing white spaces

}

func removeWhitespaces(input string) string {
	return strings.Join(strings.Fields(input), "")
}

func splitNameFromJD(linkedData *models.LinkedIn) {
	fmt.Println("data", linkedData.CompanyName)
	delimiter := " · "
	parts := strings.SplitN(linkedData.CompanyName, delimiter, 2)
	if len(parts) == 2 {
		companyName := strings.TrimSpace(parts[0])
		locationPart := strings.TrimSpace(parts[1])
		pattern := `[A-Za-z, ]+`
		re := regexp.MustCompile(pattern)

		// Find the first matching location in the string
		match := re.FindString(locationPart)
		fmt.Println("Company", companyName)
		// Trim any extra whitespace
		location := match
		location = strings.TrimSpace(location)
		linkedData.CompanyName = companyName
		linkedData.Location = location
	} else {
		fmt.Println("Invalid input format")
	}
}

// The function is responsible for sending JD to extractor engine and filter out the skills from the JD
func (c *LinkedInController) sendSkillsToEndpoint(skills string) (string, error) {
	payload := []byte(`{"text": "` + skills + `"}`)
	fmt.Println("payload", bytes.NewBuffer(payload))
	req, err := http.NewRequest("POST", c.SkillEndpoint, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	// Handle the response or error
	if err != nil {
		c.Log.WithError(err).Error("Failed to send skills data")
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Log.Errorf("Received non-OK response: %d", resp.StatusCode)
		return "", fmt.Errorf("Received non-OK response: %d", resp.StatusCode)
	}

	// Read and return the new skills from the response
	newSkills, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Log.WithError(err).Error("Failed to read response body")
		return "", err
	}

	return string(newSkills), nil
}

// Implement other CRUD operations in a similar manner.
