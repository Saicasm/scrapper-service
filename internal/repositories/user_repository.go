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

type UserRepository struct {
	Collection *mongo.Collection
	Log        *logrus.Logger
}

func NewUserRepository(collection *mongo.Collection, log *logrus.Logger) UserRepository {
	return UserRepository{
		Collection: collection,
		Log:        log,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		r.Log.WithError(err).Error("Failed to create a new user")
	}
	return err
}

func (r *UserRepository) Update(ctx context.Context, filter interface{}, update interface{}) (error, map[string]interface{}) {
	res, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.Log.WithError(err).Error("Failed to update User")
	}
	result := map[string]interface{}{"data": res}
	if res.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
	}
	if res.MatchedCount < 1 {
		r.Log.WithError(err).Error("No user Found")
	} else {
		fmt.Printf("present")

	}

	return err, result
}
func (r *UserRepository) Delete(ctx context.Context, user *models.User) error {
	_, err := r.Collection.DeleteOne(ctx, user)
	if err != nil {
		r.Log.WithError(err).Error("Failed to delete user")
	}
	return err
}

func (r *UserRepository) GetSkillsForUser(ctx context.Context, filter interface{}) (error, []string) {
	opts := options.Find().SetProjection(bson.D{{"skills", 1}, {"email", 1}, {"first_name", 1}})
	res, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		r.Log.WithError(err).Error("Failed to Get skills for user")
	}
	var results []models.User
	for res.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.User
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

	return err, results[0].Skills
}

func (r *UserRepository) GetUserById(ctx context.Context, filter interface{}) (error, []models.User) {
	opts := options.Find()
	res, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		r.Log.WithError(err).Error("Failed to Get user")
	}
	var results []models.User
	for res.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.User
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

func (r *UserRepository) GetAllUsers(ctx context.Context) (error, []models.User) {

	findOptions := options.Find()
	//Set the limit of the number of record to find
	findOptions.SetLimit(5)
	res, err := r.Collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		r.Log.WithError(err).Error("Failed to get all the users")
	}
	var results []models.User

	//Finding multiple documents returns a cursor
	//Iterate through the cursor allows us to decode documents one at a time

	for res.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.User
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
