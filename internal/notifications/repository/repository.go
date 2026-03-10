package repository

import "github.com/DuvanRozoParra/sicou/internal/notifications/models"

type Repository interface {
	Save(notification models.Notification) error
}

type InMemoryRepository struct{}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{}
}

func (r *InMemoryRepository) Save(notification models.Notification) error {
	// Simulación
	return nil
}
