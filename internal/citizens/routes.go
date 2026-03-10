package citizens

import (
	"github.com/gofiber/fiber/v2"
)

func (m *Module) RegisterRoutes(app *fiber.App) {
	citizens := app.Group("/api/v1/citizens")

	citizens.Post("/", m.HandlerCitizens.Create)
	citizens.Get("/", m.HandlerCitizens.List)
	citizens.Get("/:id", m.HandlerCitizens.GetById)
	citizens.Put("/:id", m.HandlerCitizens.Update)
	citizens.Delete("/:id", m.HandlerCitizens.Delete)
}
