package citizens

import (
	"github.com/DuvanRozoParra/sicou/internal/citizens/handlers"
	"github.com/DuvanRozoParra/sicou/internal/citizens/repository"
	"github.com/DuvanRozoParra/sicou/internal/citizens/service"
	"github.com/DuvanRozoParra/sicou/pkg/database"
)

type Module struct {
	HandlerCitizens *handlers.SCitizenHandler
}

func NewModule() *Module {
	repo_citizens := repository.NewCitizensRepository(database.DB)

	svc_citizens := service.NewCitizenService(repo_citizens)

	handler_citizens := handlers.NewCitizenHandler(svc_citizens)

	return &Module{
		HandlerCitizens: handler_citizens,
	}
}
