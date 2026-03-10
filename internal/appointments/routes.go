package appointments

import "github.com/gofiber/fiber/v2"

func (m *Module) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/v1")

	// FORM
	group.Post("/form/legal-advise", m.Handler.Create)
	group.Get("/form/legal-advise", m.Handler.List)
	group.Get("/form/legal-advise/preturnos", m.Handler.List)
	group.Get("/form/legal-advise/preturnos/assigned", m.Handler.ListAssigned)
	group.Get("/form/legal-advise/preturnos/assignment-options", m.Handler.AssignmentOptions)
	group.Patch("/form/legal-advise/preturnos/:preturno_id/assignment", m.Handler.AssignPreturno)

	// TURNOS
	turn := group.Group("/turnos")
	turn.Post("/", m.handlerTurn.Create)
	turn.Get("/", m.handlerTurn.List)
	turn.Patch("/:id", m.handlerTurn.Update)
	turn.Get("/:id", m.handlerTurn.GetById)
	turn.Delete("/:id", m.handlerTurn.Delete)
	turn.Patch("/cambio-escritorio/:id", m.handlerTurn.SetTurnDesktop)
	turn.Get("/:id/timeline", m.HandlerTurnAudit.GetTimeLine)

	// PRE TURNOS
	// pre_turn := group.Group("/pre-turn")
	//pre_turn.Get("/", m.handlerPreTurn.GetServiceTypes)
	//pre_turn.Get("/area-legal", m.handlerPreTurn.GetAreaLaboral)

	// ESCRITORIOS
	escritorios := group.Group("/escritorios")
	escritorios.Post("/", m.HandlerDesktop.Create)
	escritorios.Patch("/:code", m.HandlerDesktop.UpdateByCode)
	escritorios.Get("/", m.HandlerDesktop.List)
	escritorios.Delete("/:code", m.HandlerDesktop.DeleteByCode)
}
