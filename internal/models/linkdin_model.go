package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Status string

const (
	APPLIED     string = "APPLIED"
	IN_PROGRESS string = "IN_PROGRESS"
	REJECTED    string = "REJECTED"
	INTERESTED  string = "INTERESTED"
)

type LinkedIn struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	CompanyName    string             `json:"company_name" Usage:"required" `
	Compensation   string             `json:"compensation" Usage:"required,alphanumeric"`
	Title          string             `json:"title"`
	Skills         []string           `json:"skills"`
	JobDescription string             `json:"job_description" bson:"job_description"`
	Location       string             `json:"location"`
	Score          string             `json:"score" bson:"score"`
	UserId         string             `json:"user_id" bson:"user_id"`
	Status         Status             `json:"status"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

func NewLinkedIn() LinkedIn {
	return LinkedIn{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
