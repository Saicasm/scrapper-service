package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func IngestData(c *gin.Context) {
	//var IngestData models.Ingest
	//if err := c.ShouldBindJSON(&IngestData); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//newData, err := storage.CreateIngestData(IngestData)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	c.JSON(http.StatusCreated, gin.H{"status": "OK"})
}
