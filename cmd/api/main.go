package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/appointments"
	"github.com/DuvanRozoParra/sicou/internal/auth"
	cases "github.com/DuvanRozoParra/sicou/internal/case"
	"github.com/DuvanRozoParra/sicou/internal/citizens"
	"github.com/DuvanRozoParra/sicou/internal/dashboard"
	"github.com/DuvanRozoParra/sicou/internal/documents"
	"github.com/DuvanRozoParra/sicou/internal/pqrs"
	"github.com/DuvanRozoParra/sicou/pkg/config"
	"github.com/DuvanRozoParra/sicou/pkg/database"
	"github.com/DuvanRozoParra/sicou/pkg/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	ctx := context.Background()

	/*
		eventBus := messaging.NewKafkaEventBus(
			[]string{"localhost:9092"},
		)
		defer eventBus.Close()

	*/
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	_, err = database.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	app := fiber.New()

	allowOrigins := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS"))
	if allowOrigins == "" {
		allowOrigins = "http://localhost:5173,http://localhost:4173"
	}
	if strings.Contains(allowOrigins, "*") {
		log.Printf("CORS_ALLOW_ORIGINS cannot include wildcard '*' when credentials are enabled. Using safe default origins.")
		allowOrigins = "http://localhost:5173,http://localhost:4173"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	app.Use(recover.New())

	module_auth, err := auth.NewModule(ctx, auth.ModuleConfig{
		JWTSecret:             cfg.Auth.JWTSecret,
		TokenTTL:              cfg.Auth.TokenTTL,
		PasswordResetTokenTTL: cfg.Auth.PasswordResetTokenTTL,
		PasswordResetURL:      cfg.Auth.PasswordResetURL,
		CookieName:            cfg.Auth.CookieName,
		CookiePath:            cfg.Auth.CookiePath,
		CookieDomain:          cfg.Auth.CookieDomain,
		CookieSecure:          cfg.Auth.CookieSecure,
		CookieSameSite:        cfg.Auth.CookieSameSite,
		SuperAdminEmail:       cfg.Auth.SuperAdminEmail,
		SuperAdminDisplayName: cfg.Auth.SuperAdminDisplayName,
		SuperAdminPassword:    cfg.Auth.SuperAdminPassword,
	})
	if err != nil {
		log.Fatalf("Failed to initialize auth module: %v", err)
	}

	sessionGuard := module_auth.RequireSession()
	optionalSessionGuard := module_auth.OptionalSession()
	app.Use(func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodOptions {
			return c.Next()
		}

		if isPublicRoute(c.Method(), c.Path()) {
			fmt.Printf("%s \n", c.Path())
			return optionalSessionGuard(c)
		}

		return sessionGuard(c)
	})

	// moduleExample := notifications.NewModule(eventBus)
	// moduleExample.RegisterRoutes(app)

	minioClient := storage.NewMinioClient()
	documentsBucket := strings.TrimSpace(os.Getenv("DOCUMENTS_BUCKET"))
	if documentsBucket == "" {
		documentsBucket = "documents"
	}
	if err := storage.EnsureBucket(minioClient, documentsBucket); err != nil {
		log.Fatalf("Failed to ensure documents bucket %q: %v", documentsBucket, err)
	}
	pqrsBucket := strings.TrimSpace(os.Getenv("PQRS_BUCKET"))
	if pqrsBucket == "" {
		pqrsBucket = "pqrsd"
	}
	if err := storage.EnsureBucket(minioClient, pqrsBucket); err != nil {
		log.Fatalf("Failed to ensure pqrs bucket %q: %v", pqrsBucket, err)
	}

	appointmentsModule := appointments.NewModule(database.DB, minioClient, documentsBucket)
	appointmentsModule.RegisterRoutes(app)

	pqrsModule := pqrs.NewModule(database.DB, minioClient, pqrsBucket)
	pqrsModule.RegisterRoutes(app)
	module_documents := documents.NewModule()
	module_documents.RegisterRoutes(app)

	module_auth.RegisterRoutes(app)

	module_citizens := citizens.NewModule()
	module_citizens.RegisterRoutes(app)

	module_case := cases.NewModule()
	module_case.RegisterRoutes(app)

	module_dashboard := dashboard.NewModule()
	module_dashboard.RegisterRoutes(app)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}

	/*
		go func() {
		}()

		// Graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		_ = app.Shutdown()
	*/
}

func isPublicRoute(method string, path string) bool {
	switch {
	case method == fiber.MethodPost && (path == "/api/v1/case-file/verify-otp"):
		return true
	case method == fiber.MethodPost && (path == "/auth/login" || path == "/api/v1/auth/login"):
		return true
	case method == fiber.MethodPost && (path == "/auth/forgot-password" || path == "/api/v1/auth/forgot-password"):
		return true
	case method == fiber.MethodPost && (path == "/auth/reset-password" || path == "/api/v1/auth/reset-password"):
		return true
	case method == fiber.MethodPost && (path == "/auth/logout" || path == "/api/v1/auth/logout"):
		return true
	case method == fiber.MethodPost && path == "/api/v1/form/legal-advise":
		return true
	default:
		return false
	}
}
