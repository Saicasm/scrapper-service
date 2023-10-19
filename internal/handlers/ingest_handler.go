package handlers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// No Concurrency design
//func IngestData(c *gin.Context) {
//
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
//		return
//	}
//	htmlContent := string(body)
//	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse HTML content"})
//		return
//	}
//
//	// Extract data from elements with a specific CSS class
//	extractedData := []string{}
//	doc.Find(".job-details-jobs-unified-top-card__job-title").Each(func(i int, s *goquery.Selection) {
//		// Extract and add the text or attributes you need to the result slice
//		text := s.Text()
//		extractedData = append(extractedData, text)
//	})
//
//	doc.Find(".job-details-jobs-unified-top-card__primary-description").Each(func(i int, s *goquery.Selection) {
//		// Extract and add the text or attributes you need to the result slice
//		text := s.Text()
//		extractedData = append(extractedData, text)
//	})
//	doc.Find(".job-details-how-you-match__skills-item-subtitle").Each(func(i int, s *goquery.Selection) {
//		// Extract and add the text or attributes you need to the result slice
//		text := s.Text()
//		extractedData = append(extractedData, text)
//	})
//	doc.Find(".jobs-description-content__text--stretch").Each(func(i int, s *goquery.Selection) {
//		// Extract and add the text or attributes you need to the result slice
//		text := s.Text()
//		extractedData = append(extractedData, text)
//	})
//
//	fmt.Println("Title:", extractedData)
//	c.JSON(http.StatusCreated, gin.H{"status": "OK"})
//}

func IngestData(c *gin.Context) {

	// Read the request body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	htmlContent := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse HTML content"})
		return
	}

	// Create a channel to receive extracted data
	dataChan := make(chan string)

	// Launch Goroutines to extract data from different elements
	go extractData(doc, ".job-details-jobs-unified-top-card__job-title", dataChan)
	go extractData(doc, ".job-details-jobs-unified-top-card__primary-description", dataChan)
	go extractData(doc, ".job-details-how-you-match__skills-item-subtitle", dataChan)
	go extractData(doc, ".jobs-description-content__text--stretch", dataChan)

	var extractedData []string

	// Collect data from the channel
	for i := 0; i < 4; i++ { // Adjust this number based on the number of Goroutines launched
		data := <-dataChan
		extractedData = append(extractedData, data)

	}

	// Close the channel
	//close(dataChan)

	fmt.Println("Title:", extractedData)
	c.JSON(http.StatusCreated, gin.H{"status": "OK"})
}

func extractData(doc *goquery.Document, selector string, dataChan chan string) {
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		dataChan <- text
	})
}
