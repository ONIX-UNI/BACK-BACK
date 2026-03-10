package dashboard

import "github.com/gofiber/fiber/v2"

func (m *Module) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/v1/dashboard")
	group.Get("/overview", m.Handler.Overview)
}
