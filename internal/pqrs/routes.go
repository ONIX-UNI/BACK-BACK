package pqrs

import "github.com/gofiber/fiber/v2"

func (m *Module) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/v1")
	group.Post("/pqrs", m.Handler.Create)
	group.Get("/pqrs", m.Handler.List)
}
