package services

import (
	"context"
	"github.com/scraper/internal/models"
	"github.com/scraper/internal/repositories"
	"github.com/sirupsen/logrus"
	"time"
)

type LinkedInService struct {
	Repository repositories.LinkedInRepository
	Log        *logrus.Logger
}

func NewLinkedInService(repository repositories.LinkedInRepository, log *logrus.Logger) LinkedInService {
	return LinkedInService{
		Repository: repository,
		Log:        log,
	}
}

func (s *LinkedInService) Create(linkedin *models.LinkedIn) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Repository.Create(ctx, linkedin)
	if err != nil {
		s.Log.WithError(err).Error("Failed to create a new todo")
	}
	return err
}
