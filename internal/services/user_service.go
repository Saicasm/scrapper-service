package services

import (
	"context"
	"github.com/scraper/internal/models"
	"github.com/scraper/internal/repositories"
	"github.com/sirupsen/logrus"
	"time"
)

type UserService struct {
	Repository repositories.UserRepository
	Log        *logrus.Logger
}

func NewUserService(repository repositories.UserRepository, log *logrus.Logger) UserService {
	return UserService{
		Repository: repository,
		Log:        log,
	}
}

func (s *UserService) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Repository.Create(ctx, user)
	if err != nil {
		s.Log.WithError(err).Error("Failed to create new user")
	}
	return err
}
func (s *UserService) Update(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Repository.Create(ctx, user)
	if err != nil {
		s.Log.WithError(err).Error("Failed to update user")
	}
	return err
}
