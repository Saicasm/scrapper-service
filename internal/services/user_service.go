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
func (s *UserService) Update(filter interface{}, update interface{}) (error, map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err, res := s.Repository.Update(ctx, filter, update)
	if err != nil {
		s.Log.WithError(err).Error("Failed to update user")
	}
	return err, res
}

func (s *UserService) GetAllUsers() (error, []models.User) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err, result := s.Repository.GetAllUsers(ctx)
	if err != nil {
		s.Log.WithError(err).Error("Failed to get all users")
	}
	return err, result
}

func (s *UserService) GetUserSkills(filter interface{}) (error, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err, result := s.Repository.GetSkillsForUser(ctx, filter)
	if err != nil {
		s.Log.WithError(err).Error("Failed to get all users")
	}
	return err, result
}
func (s *UserService) GetUserById(filter interface{}) (error, []models.User) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err, result := s.Repository.GetUserById(ctx, filter)
	if err != nil {
		s.Log.WithError(err).Error("Failed to get all users")
	}
	return err, result
}
