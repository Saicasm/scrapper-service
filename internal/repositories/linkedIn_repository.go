package repositories

import (
	"context"
	"github.com/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkedInRepository struct {
	Collection *mongo.Collection
	Log        *logrus.Logger
}

func NewLinkedInRepository(collection *mongo.Collection, log *logrus.Logger) LinkedInRepository {
	return LinkedInRepository{
		Collection: collection,
		Log:        log,
	}
}

func (r *LinkedInRepository) Create(ctx context.Context, linkedin *models.LinkedIn) error {
	_, err := r.Collection.InsertOne(ctx, linkedin)
	if err != nil {
		r.Log.WithError(err).Error("Failed to create a new job")
	}
	return err
}

// Implement other CRUD operations (Get, Update, Delete) in a similar manner.
