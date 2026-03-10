package notifications

import "github.com/gofiber/fiber/v2"

func (m *Module) RegisterRoutes(app *fiber.App) {
	group := app.Group("/notifications")

	group.Post("/hello", m.Handler.Hello)
}
