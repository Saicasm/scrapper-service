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

func IngestData(c *gin.Context) {

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

	// Extract data from elements with a specific CSS class
	extractedData := []string{}
	doc.Find(".job-details-jobs-unified-top-card__job-title").Each(func(i int, s *goquery.Selection) {
		// Extract and add the text or attributes you need to the result slice
		text := s.Text()
		extractedData = append(extractedData, text)
	})

	doc.Find(".job-details-jobs-unified-top-card__primary-description").Each(func(i int, s *goquery.Selection) {
		// Extract and add the text or attributes you need to the result slice
		text := s.Text()
		extractedData = append(extractedData, text)
	})
	doc.Find(".job-details-how-you-match__skills-item-subtitle").Each(func(i int, s *goquery.Selection) {
		// Extract and add the text or attributes you need to the result slice
		text := s.Text()
		extractedData = append(extractedData, text)
	})
	doc.Find(".jobs-description-content__text--stretch").Each(func(i int, s *goquery.Selection) {
		// Extract and add the text or attributes you need to the result slice
		text := s.Text()
		extractedData = append(extractedData, text)
	})

	fmt.Println("Title:", extractedData)
	c.JSON(http.StatusCreated, gin.H{"status": "OK"})
}
