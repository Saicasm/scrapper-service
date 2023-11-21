package models

import (
	"github.com/google/uuid"
	"time"
)

type LinkedIn struct {
	ID             uuid.UUID `json:"_id" bson:"_id"`
	CompanyName    string    `json:"company_name" Usage:"required"`
	Compensation   string    `json:"compensation" Usage:"required,alphanumeric"`
	Title          string    `json:"title"`
	Skills         string    `json:"skills"`
	JobDescription string    `json:"job_description"`
	Location       string    `json:"location"`
	Score          string    `json:"score"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewLinkedIn() LinkedIn {
	return LinkedIn{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
