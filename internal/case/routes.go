package cases

import "github.com/gofiber/fiber/v2"

func (m *Module) RegisterRoutes(app *fiber.App) {
	file := app.Group("/api/v1/case-file")

	folder := app.Group("/api/v1/folders")

	folder.Get("/", m.HandlerCaseFile.FolderList)

	file.Post("/", m.HandlerCaseFile.Create)
	file.Put("/:id", m.HandlerCaseFile.Update)
	file.Get("/list", m.HandlerCaseFile.ListCase)
	file.Get("/", m.HandlerCaseFile.List)
	file.Get("/by-email", m.HandlerCaseFile.ListCasesByEmail)
	file.Get("/:id", m.HandlerCaseFile.GetById)
	file.Delete("/:id", m.HandlerCaseFile.Delete)

	file.Post("/request-otp", m.HandlerCaseFile.CaseOTP)
	file.Post("/verify-otp", m.HandlerCaseFile.VerifyOTP)

	events := app.Group("/case-events")

	events.Post("/", m.HandlerCaseEvent.Create)
	events.Get("/", m.HandlerCaseEvent.List)
	events.Get("/:id", m.HandlerCaseEvent.GetById)
	events.Put("/:id", m.HandlerCaseEvent.Update)
	events.Delete("/:id", m.HandlerCaseEvent.Delete)

	checkListItemDef := app.Group("/check-list")

	checkListItemDef.Post("/", m.handlerCheckListDef.Create)
	checkListItemDef.Get("/", m.handlerCheckListDef.List)
	checkListItemDef.Get("/:id", m.handlerCheckListDef.GetById)
	checkListItemDef.Put("/:id", m.handlerCheckListDef.Update)
	checkListItemDef.Delete("/:id", m.handlerCheckListDef.Delete)

	serviceType := app.Group("/service-type")

	serviceType.Post("/", m.handlerServiceType.Create)
	serviceType.Get("/", m.handlerServiceType.List)
	serviceType.Get("/:id", m.handlerServiceType.GetById)
	serviceType.Put("/:id", m.handlerServiceType.Update)
	serviceType.Delete("/:id", m.handlerServiceType.Delete)

	legalArea := app.Group("/legal-area")

	legalArea.Post("/", m.handlerLegalArea.Create)
	legalArea.Get("/", m.handlerLegalArea.List)
	legalArea.Get("/:id", m.handlerLegalArea.GetById)
	legalArea.Put("/:id", m.handlerLegalArea.Update)
	legalArea.Delete("/:id", m.handlerLegalArea.Delete)
}
