package service

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/notifications/events"
	"github.com/DuvanRozoParra/sicou/internal/notifications/models"
	"github.com/DuvanRozoParra/sicou/internal/notifications/repository"
	"github.com/DuvanRozoParra/sicou/pkg/messaging"
	"github.com/google/uuid"
)

type Service struct {
	repo      repository.Repository
	publisher messaging.Publisher
}

func NewService(
	repo repository.Repository,
	publisher messaging.Publisher,
) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) SendHello(ctx context.Context, name string) error {

	notification := models.Notification{
		ID:   uuid.NewString(),
		Name: name,
		Type: "HELLO_EVENT",
	}

	if err := s.repo.Save(notification); err != nil {
		return err
	}

	event := events.HelloEvent{
		ID:   notification.ID,
		Name: notification.Name,
		Type: notification.Type,
	}

	return s.publisher.Publish(ctx, event)
}
