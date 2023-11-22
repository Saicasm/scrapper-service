package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `json:"first_name" bson:"first_name" Usage:"required"`
	LastName  string             `json:"last_name" bson:"last_name" Usage:"required"`
	Email     string             `json:"email" bson:"email"`
	Skills    []string           `json:"skills" bson:"skills"`
	Location  string             `json:"location" bson:"location"`
	CreatedAt time.Time          `json:"created_at" `
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewUser() User {
	return User{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
