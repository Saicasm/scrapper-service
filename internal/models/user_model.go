package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"_id" bson:"_id"`
	FirstName string    `json:"first_name" Usage:"required"`
	LastName  string    `json:"last_name" Usage:"required"`
	Skills    []string  `json:"skills"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser() User {
	return User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
