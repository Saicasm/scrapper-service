package repositories

import (
	"context"
	"fmt"
	"github.com/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
func (r *LinkedInRepository) Update(ctx context.Context, filter interface{}, update interface{}) (error, map[string]interface{}) {
	res, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.Log.WithError(err).Error("Failed to update Job")
	}
	result := map[string]interface{}{"data": res}
	if res.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
	}
	if res.MatchedCount < 1 {
		r.Log.WithError(err).Error("No Job Found")
	} else {
	}

	return err, result
}

func (r *LinkedInRepository) GetJobsForUser(ctx context.Context, filter interface{}) (error, []models.LinkedIn) {
	findOptions := options.Find().SetProjection(bson.D{{"job_description", 0}}).SetLimit(100)

	//Set the limit of the number of record to find
	//findOptions.SetLimit(100)
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

func (r *LinkedInRepository) GetAnalyticsForUser(ctx context.Context, filter interface{}) (error, interface{}) {
	filter1 := bson.D{{"user_id", bson.D{{"$eq", filter}}}}
	countInterested := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user_id", bson.D{{"$eq", filter}}}},
				bson.D{{"status", bson.D{{"$eq", models.INTERESTED}}}},
			},
		},
	}
	countInProgress := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user_id", bson.D{{"$eq", filter}}}},
				bson.D{{"status", bson.D{{"$eq", models.IN_PROGRESS}}}},
			},
		},
	}
	countRejected := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user_id", bson.D{{"$eq", filter}}}},
				bson.D{{"status", bson.D{{"$eq", models.REJECTED}}}},
			},
		},
	}
	countApplied := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user_id", bson.D{{"$eq", filter}}}},
				bson.D{{"status", bson.D{{"$eq", models.APPLIED}}}},
			},
		},
	}
	count, err := r.Collection.CountDocuments(context.TODO(), filter1)
	if err != nil {
		panic(err)
	}
	countResApplied, err := r.Collection.CountDocuments(context.TODO(), countApplied)
	if err != nil {
		panic(err)
	}
	countResRejected, err := r.Collection.CountDocuments(context.TODO(), countRejected)
	if err != nil {
		panic(err)
	}
	countResInProgress, err := r.Collection.CountDocuments(context.TODO(), countInProgress)
	if err != nil {
		panic(err)
	}
	countResIntrested, err := r.Collection.CountDocuments(context.TODO(), countInterested)
	if err != nil {
		panic(err)
	}
	result := map[string]interface{}{"count": count, "applied": countResApplied, "rejected": countResRejected, "in_progress": countResInProgress, "interested": countResIntrested}

	return err, result
}
