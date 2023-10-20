package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoDB *mongo.Database

func InitMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // MongoDB connection string
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	MongoDB = client.Database("jts") // Replace "todo_app" with your database name
}
