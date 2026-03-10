package appointments

import (
	"github.com/DuvanRozoParra/sicou/internal/appointments/handlers"
	"github.com/DuvanRozoParra/sicou/internal/appointments/repository"
	"github.com/DuvanRozoParra/sicou/internal/appointments/service"
	documentsrepository "github.com/DuvanRozoParra/sicou/internal/documents/repository"
	documentservice "github.com/DuvanRozoParra/sicou/internal/documents/service"
	"github.com/DuvanRozoParra/sicou/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Handler          *handlers.Handler
	handlerTurn      *handlers.STurnHandler
	HandlerTurnAudit *handlers.STurnAuditHandler
	HandlerDesktop   *handlers.SDesktopHandler
	// handlerPreTurn   *handlers.SPreturnHandlers
}

func NewModule(db *pgxpool.Pool, minioClient *storage.Client, bucket string) *Module {

	repo := repository.NewPostgresRepository(db)
	repo_turn := repository.NewTurnRepository(db)
	repo_turn_audit := repository.NewTurnAuditRepository(db)
	//repo_pre_trun := repository.NewPreturnRepository(db)
	repo_desktop := repository.NewDesktopRepository(db)

	fileRepo := documentsrepository.NewPostgresRepository(db)
	fileSvc := documentservice.NewFileObjectService(fileRepo, minioClient, bucket)

	svc := service.NewService(repo, fileSvc)
	svc_turn := service.NewTurnService(repo_turn)
	svc_turn_audit := service.NewTurnAuditService(repo_turn_audit)
	svc_desktop := service.NewDesktopService(repo_desktop)
	//svc_pre_turn := service.NewPreTurnService(repo_pre_trun)

	handler := handlers.NewHandler(svc)
	handler_turn := handlers.NewTurnHandler(svc_turn)
	handler_turn_audit := handlers.NewTurnAuditHandler(svc_turn_audit)
	handler_desktop := handlers.NewDesktopHandler(svc_desktop)
	//handler_pre_turn := handlers.NewPreTurnHandler(svc_pre_turn)

	return &Module{
		Handler:          handler,
		handlerTurn:      handler_turn,
		HandlerTurnAudit: handler_turn_audit,
		HandlerDesktop:   handler_desktop,
		// handlerPreTurn:   handler_pre_turn,
	}
}
