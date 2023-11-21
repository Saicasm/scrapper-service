// server/controllers/todo_controller.go
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/scraper/internal/models"
	"github.com/scraper/internal/services"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type UserController struct {
	Service       services.UserService
	Log           *logrus.Logger
	SkillEndpoint string // Define the endpoint URL for skills
}

func NewUserController(service services.UserService, log *logrus.Logger, skillEndpoint string) *UserController {
	return &UserController{
		Service:       service,
		Log:           log,
		SkillEndpoint: skillEndpoint, // Initialize the endpoint URL
	}
}

func (c *UserController) UpdateUser(ctx *gin.Context) {

}
func (c *UserController) Create(ctx *gin.Context) {
	c.Log.Debug("Create User")
	user := models.NewUser()
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}
	if err := c.Service.Create(&user); err != nil {
		c.Log.WithError(err).Error("Failed to create a new job")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new job"})
	} else {
		ctx.JSON(http.StatusCreated, user)
	}
}
