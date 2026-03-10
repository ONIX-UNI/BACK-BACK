package documents

import (
	"github.com/DuvanRozoParra/sicou/internal/documents/handlers"
	"github.com/DuvanRozoParra/sicou/internal/documents/repository"
	services "github.com/DuvanRozoParra/sicou/internal/documents/service"
	"github.com/DuvanRozoParra/sicou/pkg/database"
	"github.com/DuvanRozoParra/sicou/pkg/storage"
)

type Module struct {
	HandlerDocFileObject *handlers.SDocFileObjectHandler
	HandlerFileObject    *handlers.FileObjectHandler
	HandlerDocs          *handlers.DocumentHandler
}

func NewModule() *Module {
	// 🔹 MinIO
	minioClient := storage.NewMinioClient()
	bucket := "documents"
	storage.EnsureBucket(minioClient, bucket)
	bucketExpedient := "expedient-docs"
	storage.EnsureBucket(minioClient, bucketExpedient)

	repo_file_object := repository.NewPostgresRepository(database.DB)
	repo_docs := repository.NewPostgresDocumentRepository(database.DB)
	repo_docs_file_object := repository.NewDocFileObjectRepository(database.DB)

	svc_file_object := services.NewFileObjectService(
		repo_file_object,
		minioClient,
		bucket,
	)
	svc_docs_file_object := services.NewDocFileObjectService(repo_docs_file_object, minioClient, bucketExpedient)
	svc_docs := services.NewDocumentService(repo_docs)

	handler_file_object := handlers.NewFileObjectHandler(svc_file_object)
	handler_docs := handlers.NewDocumentHandler(svc_docs)
	handler_docs_file_object := handlers.NewDocFileObjectHandler(svc_docs_file_object)

	return &Module{
		HandlerDocs:          handler_docs,
		HandlerFileObject:    handler_file_object,
		HandlerDocFileObject: handler_docs_file_object,
	}
}
