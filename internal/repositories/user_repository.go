package repositories

import (
	"context"
	"github.com/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
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

//	func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
//		_, err := r.Collection.UpdateOne(ctx, user)
//		if err != nil {
//			r.Log.WithError(err).Error("Failed to create a new job")
//		}
//		return err
//	}
func (r *UserRepository) Delete(ctx context.Context, user *models.User) error {
	_, err := r.Collection.DeleteOne(ctx, user)
	if err != nil {
		r.Log.WithError(err).Error("Failed to delete user")
	}
	return err
}

func (r *UserRepository) GetAllUsers(ctx context.Context, user *models.User) error {
	_, err := r.Collection.Find(ctx, user)
	if err != nil {
		r.Log.WithError(err).Error("Failed to get all the users")
	}
	return err
}
