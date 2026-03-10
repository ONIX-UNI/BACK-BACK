package dashboard

import (
	"github.com/DuvanRozoParra/sicou/internal/dashboard/handlers"
	"github.com/DuvanRozoParra/sicou/internal/dashboard/repository"
	"github.com/DuvanRozoParra/sicou/internal/dashboard/service"
	"github.com/DuvanRozoParra/sicou/pkg/database"
)

type Module struct {
	Handler *handlers.Handler
}

func NewModule() *Module {
	repo := repository.NewRepository(database.DB)
	svc := service.NewService(repo)
	handler := handlers.NewHandler(svc)

	return &Module{Handler: handler}
}
