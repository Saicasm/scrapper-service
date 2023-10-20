// internal/mongodb/mongodb.go
package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient initializes a MongoDB client and returns a client instance.
func NewClient(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}

// GetCollection returns a MongoDB collection based on the provided database and collection names.
func GetCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	return collection
}
