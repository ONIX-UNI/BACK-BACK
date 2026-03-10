package auth

import "github.com/gofiber/fiber/v2"

func (m *Module) RegisterRoutes(app *fiber.App) {
	// auth/session
	m.registerSessionRoutes(app.Group("/auth"))
	m.registerSessionRoutes(app.Group("/api/v1/auth"))

	// users
	m.registerUserRoutes(app.Group("/user"))
	m.registerUserRoutes(app.Group("/api/v1/user"))
}

func (m *Module) registerSessionRoutes(auth fiber.Router) {
	auth.Post("/login", m.HandlerAppUser.Login)
	auth.Post("/forgot-password", m.HandlerAppUser.ForgotPassword)
	auth.Post("/reset-password", m.HandlerAppUser.ResetPassword)
	auth.Get("/me", m.HandlerAppUser.Me)
	auth.Post("/refresh", m.HandlerAppUser.Refresh)
	auth.Post("/logout", m.HandlerAppUser.Logout)
}

func (m *Module) registerUserRoutes(user fiber.Router) {
	user.Use(m.HandlerAppUser.RequireUserManagementRole)

	user.Post("/", m.HandlerAppUser.Create)
	user.Post("/bulk", m.HandlerAppUser.BulkCreate)
	user.Get("/audit-log", m.HandlerAppUser.ListAuditLog)
	user.Get("/", m.HandlerAppUser.List)
	user.Get("/search", m.HandlerAppUser.GetByDisplayName)
	user.Get("/options", m.HandlerAppUser.ListOptions)
	user.Get("/email/:email", m.HandlerAppUser.GetByEmail)
	user.Get("/:id", m.HandlerAppUser.GetByID)
	user.Put("/:id", m.HandlerAppUser.Update)
	user.Delete("/:id", m.HandlerAppUser.Delete)
}
