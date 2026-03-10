package documents

import "github.com/gofiber/fiber/v2"

func (m *Module) RegisterRoutes(app *fiber.App) {
	registerRoutes := func(prefix string) {
		docs := app.Group(prefix)

		// intermediente
		doc_file_object := docs.Group("/expedient")

		doc_file_object.Post("/", m.HandlerDocFileObject.UploadDoc)
		doc_file_object.Post("/file-folder", m.HandlerDocFileObject.UploadDocFolder)

		doc_file_object.Patch("/:document_id", m.HandlerDocFileObject.UpdateDocument)
		doc_file_object.Delete("/:document_id", m.HandlerDocFileObject.DeleteDocument)
		doc_file_object.Get("/:case_id/timeline", m.HandlerDocFileObject.GetCaseTimeLineItem)
		doc_file_object.Get("/:case_id/files", m.HandlerDocFileObject.ListCaseFiles)
		doc_file_object.Get("/:file_id/view", m.HandlerDocFileObject.ViewFile)

		// DOCUMENT KIND GET
		doc_kind := docs.Group("/doc-kind")
		doc_kind.Get("/", m.HandlerDocFileObject.GetDocumentKinds)

		// File Objects
		// files.Put("/:id", m.HandlerFileObject.Update)
		files := docs.Group("/files")
		files.Post("/", m.HandlerFileObject.Create)
		files.Get("/", m.HandlerFileObject.GetAll)
		files.Get("/list", m.HandlerFileObject.List)
		files.Get("/:id/view", m.HandlerFileObject.View)
		files.Get("/:id/download", m.HandlerFileObject.Download)
		files.Get("/:id/url", m.HandlerFileObject.GetPresignedURL)
		files.Get("/:id", m.HandlerFileObject.GetById)
		files.Delete("/:id", m.HandlerFileObject.Delete)

		// Documents
		documents := docs.Group("/documents")
		documents.Post("/", m.HandlerDocs.Create)
		documents.Get("/", m.HandlerDocs.Get)
		documents.Get("/:id", m.HandlerDocs.GetByID)
		documents.Delete("/:id", m.HandlerDocs.Delete)
	}

	registerRoutes("/docs")
	registerRoutes("/api/v1/documents")
}
