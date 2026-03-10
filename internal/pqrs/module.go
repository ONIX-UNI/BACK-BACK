package pqrs

import (
	documentsrepository "github.com/DuvanRozoParra/sicou/internal/documents/repository"
	documentservice "github.com/DuvanRozoParra/sicou/internal/documents/service"
	"github.com/DuvanRozoParra/sicou/internal/pqrs/handlers"
	"github.com/DuvanRozoParra/sicou/internal/pqrs/repository"
	"github.com/DuvanRozoParra/sicou/internal/pqrs/service"
	"github.com/DuvanRozoParra/sicou/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Handler *handlers.Handler
}

func NewModule(db *pgxpool.Pool, minioClient *storage.Client, bucket string) *Module {
	repo := repository.NewPostgresRepository(db)

	fileRepo := documentsrepository.NewPostgresRepository(db)
	fileSvc := documentservice.NewFileObjectService(fileRepo, minioClient, bucket)

	svc := service.NewService(repo, fileSvc)
	handler := handlers.NewHandler(svc)

	return &Module{
		Handler: handler,
	}
}
