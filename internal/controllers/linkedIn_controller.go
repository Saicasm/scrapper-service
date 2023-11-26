// server/controllers/todo_controller.go
package controllers

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/scraper/internal/models"
	"github.com/scraper/internal/services"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type LinkedInController struct {
	Service       services.LinkedInService
	Log           *logrus.Logger
	SkillEndpoint string // Define the endpoint URL for skills
}
type extractorResponse struct {
	TechnicalTerms []string `json:"technical_terms"`
	Score          string   `json:"score"`
}

const userSkillsEndpoint = "http://localhost:8080/api/v1/ingest/user/skills/saicsm@gmail.com"

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
	//skillsChan := make(chan []string)
	compensationChan := make(chan string)
	companyChan := make(chan string)

	// Launch Goroutines to extract data from different elements
	// TODO: Extract skills from the description  chan
	go extractData(doc, ".job-details-jobs-unified-top-card__job-title", titleChan)
	go extractData(doc, ".jobs-description__content", descriptionChan)
	//go extractData(doc, ".jobs-description__content", skillsChan)
	go extractData(doc, ".job-details-jobs-unified-top-card__job-insight", compensationChan)
	go extractData(doc, ".job-details-jobs-unified-top-card__primary-description", companyChan)
	// Receive extracted data from channels

	// Store the data into the structs
	linkedin.Title = <-titleChan
	linkedin.JobDescription = <-descriptionChan
	//linkedin.Skills = <-skillsChan
	linkedin.Compensation = <-compensationChan
	linkedin.CompanyName = <-companyChan
	userid := ctx.Param("userId")
	//TODO: Get the data from the request
	linkedin.UserId = userid
	linkedin.Status = models.Status(models.INTERESTED)
	// Clean the data
	cleanData(&linkedin.Title)
	cleanData(&linkedin.JobDescription)
	//cleanData(&linkedin.Skills)
	cleanData(&linkedin.Compensation)

	extractedSkills := linkedin.JobDescription
	cleanedText := removeWhitespaces(linkedin.Compensation)

	cleanedComp := extractCompensation(cleanedText)

	linkedin.Compensation = cleanedComp
	userSkills, err := c.getSkillsForUser(userid)
	fmt.Printf("%v", userSkills)
	newSkills, score, err := c.sendSkillsToEndpoint(extractedSkills, userSkills)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send skills data"})
		return
	}
	// Update the LinkedIn model with the new skills
	linkedin.Skills = newSkills
	linkedin.Score = score
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
		return cleanedText[:(compensationHyphenIndex + (compensationHyphenIndex - 1) + 3)]
	} else {
		compensationRIndex := strings.IndexByte(cleanedText, 'r')

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
	delimiter := " · "
	parts := strings.SplitN(linkedData.CompanyName, delimiter, 2)
	if len(parts) == 2 {
		companyName := strings.TrimSpace(parts[0])
		locationPart := strings.TrimSpace(parts[1])
		pattern := `[A-Za-z, ]+`
		re := regexp.MustCompile(pattern)

		// Find the first matching location in the string
		match := re.FindString(locationPart)
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
func (c *LinkedInController) sendSkillsToEndpoint(skills string, userSkills []string) ([]string, string, error) {

	data := map[string]interface{}{
		"text":       skills,
		"userSkills": userSkills,
	}

	// Convert the map to a JSON string
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
	}

	req, err := http.NewRequest("POST", c.SkillEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("resp", resp)
	// Handle the response or error
	if err != nil {
		c.Log.WithError(err).Error("Failed to send skills data")
		return make([]string, 0), string(rune(0)), err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.Log.Errorf("Received non-OK response: %d", resp.StatusCode)
		return make([]string, 0), string(rune(0)), fmt.Errorf("Received non-OK response: %d", resp.StatusCode)
	}

	// Read and return the new skills from the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Log.WithError(err).Error("Failed to read response body")
		return make([]string, 0), string(rune(0)), err
	}
	var extractorResp extractorResponse
	err = json.Unmarshal([]byte(body), &extractorResp)
	if err != nil {
		fmt.Println("Error unmarshalling score:", err)
		return make([]string, 0), string(rune(0)), err
	}
	return extractorResp.TechnicalTerms, extractorResp.Score, err
}

func (c *LinkedInController) getSkillsForUser(userId string) ([]string, error) {

	req, err := http.NewRequest("GET", userSkillsEndpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)

	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	// Handle the response or error
	if err != nil {
		c.Log.WithError(err).Error("Failed to send skills data")
		return make([]string, 0), err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Log.Errorf("Received non-OK response: %d", resp.StatusCode)
		return make([]string, 0), fmt.Errorf("Received non-OK response: %d", resp.StatusCode)
	}

	// Read and return the new skills from the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Log.WithError(err).Error("Failed to read response body")
		return make([]string, 0), err
	}
	var skills []string
	err = json.Unmarshal(body, &skills)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return nil, nil
	}
	return skills, nil
}

func (c *LinkedInController) GetJobsForUserID(ctx *gin.Context) {
	userid := ctx.Param("userId")
	filter := bson.M{"user_id": userid}
	c.Log.Debug("Get Jobs For UserID ")
	if err, result := c.Service.GetJobsForUser(filter); err != nil {
		c.Log.WithError(err).Error("Failed to create a new job")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new job"})
	} else {
		ctx.JSON(http.StatusOK, result)
	}
}
func (c *LinkedInController) GetAnalyticsForUser(ctx *gin.Context) {
	userid := ctx.Param("userId")
	c.Log.Debug("Get Jobs For UserID ")
	if err, result := c.Service.GetAnalyticsForUser(userid); err != nil {
		c.Log.WithError(err).Error("Failed to get analytics for the user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analytics for the user"})
	} else {
		ctx.JSON(http.StatusOK, result)
	}
}
func (c *LinkedInController) UpdateJob(ctx *gin.Context) {
	userid := ctx.Param("jobId")
	objId, _ := primitive.ObjectIDFromHex(userid)
	filter := bson.M{"_id": objId}
	job := models.NewLinkedIn()
	if err := ctx.BindJSON(&job); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}
	update := bson.D{{"$set", bson.D{{"status", job.Status}, {"compensation", job.Compensation}, {"user_id", job.UserId}, {"updated_at", time.Now()}}}}
	if err, res := c.Service.Update(filter, update); err != nil {
		c.Log.WithError(err).Error("Failed to create a new job")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new job"})
	} else {
		ctx.JSON(http.StatusOK, res)
	}
}
