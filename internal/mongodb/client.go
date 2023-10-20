// mongodb/client.go

package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient represents a MongoDB client.
type MongoDBClient struct {
	client *mongo.Client
}

// NewMongoDBClient creates a new MongoDB client based on the provided configuration.
func NewMongoDBClient(config *MongoDBConfig) (*MongoDBClient, error) {
	clientOptions := options.Client().ApplyURI(config.ConnectionURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return &MongoDBClient{client: client}, nil
}

// GetClient returns the MongoDB client instance.
func (c *MongoDBClient) GetClient() *mongo.Client {
	return c.client
}
