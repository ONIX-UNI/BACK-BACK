package notifications

import (
	"github.com/DuvanRozoParra/sicou/internal/notifications/handlers"
	"github.com/DuvanRozoParra/sicou/internal/notifications/repository"
	"github.com/DuvanRozoParra/sicou/internal/notifications/service"
	"github.com/DuvanRozoParra/sicou/pkg/messaging"
)

type Module struct {
	Handler *handlers.Handler
}

func NewModule(publisher messaging.Publisher) *Module {

	repo := repository.NewInMemoryRepository()

	svc := service.NewService(repo, publisher)

	handler := handlers.NewHandler(svc)

	return &Module{
		Handler: handler,
	}
}
