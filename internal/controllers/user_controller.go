// server/controllers/todo_controller.go
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/scraper/internal/models"
	"github.com/scraper/internal/services"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
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

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userid := ctx.Param("userId")
	objId, _ := primitive.ObjectIDFromHex(userid)
	filter := bson.M{"_id": objId}
	user := models.NewUser()
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}
	update := bson.D{{"$set", bson.D{{"first_name", user.FirstName}, {"last_name", user.LastName}, {"skills", user.Skills}, {"location", user.Location}, {"email", user.Email}, {"updated_at", time.Now()}}}}
	if err, res := c.Service.Update(filter, update); err != nil {
		c.Log.WithError(err).Error("Failed to create a new job")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new job"})
	} else {
		ctx.JSON(http.StatusOK, res)
	}
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	c.Log.Debug("Get All Users")
	if err, result := c.Service.GetAllUsers(); err != nil {
		c.Log.WithError(err).Error("Failed to create a new job")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new job"})
	} else {
		ctx.JSON(http.StatusOK, result)
	}
}

func (c *UserController) GetSkillsForUser(ctx *gin.Context) {
	userid := ctx.Param("userId")
	filter := bson.M{"email": userid}
	c.Log.Debug("Get All Users")
	if err, result := c.Service.GetUserSkills(filter); err != nil {
		c.Log.WithError(err).Error("Failed to create a new job")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new job"})
	} else {
		ctx.JSON(http.StatusOK, result)
	}
}
