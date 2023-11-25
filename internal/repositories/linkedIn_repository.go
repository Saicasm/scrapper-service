package repositories

import (
	"context"
	"github.com/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

func (r *LinkedInRepository) GetJobsForUser(ctx context.Context, filter interface{}) (error, []models.LinkedIn) {

	findOptions := options.Find()
	//Set the limit of the number of record to find
	findOptions.SetLimit(5)
	res, err := r.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		r.Log.WithError(err).Error("Failed to get all the users")
	}
	var results []models.LinkedIn

	for res.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.LinkedIn
		err := res.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)

	}

	if err := res.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	res.Close(context.TODO())

	return err, results
}
