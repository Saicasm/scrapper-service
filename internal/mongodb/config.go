// mongodb/config.go

package mongodb

// MongoDBConfig represents the MongoDB database configuration.
type MongoDBConfig struct {
	ConnectionURL string
	DatabaseName  string
}

// NewMongoDBConfig returns a new MongoDBConfig with default values.
func NewMongoDBConfig() *MongoDBConfig {
	return &MongoDBConfig{
		ConnectionURL: "mongodb://localhost:27017",
		DatabaseName:  "gojob",
	}
}
