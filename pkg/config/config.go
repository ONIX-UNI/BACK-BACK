package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Port   string
	SMTP   SMTPConfig
	Auth   AuthConfig
}

type SMTPConfig struct {
	Host               string
	Port               int
	Username           string
	Password           string
	AccessToken        string
	AuthMethod         string
	FromName           string
	FromAddress        string
	UseStartTLS        bool
	UseImplicitTLS     bool
	RequireTLS         bool
	InsecureSkipVerify bool
	Timeout            time.Duration
}

type AuthConfig struct {
	JWTSecret             string
	TokenTTL              time.Duration
	PasswordResetTokenTTL time.Duration
	PasswordResetURL      string
	CookieName            string
	CookiePath            string
	CookieDomain          string
	CookieSecure          bool
	CookieSameSite        string
	SuperAdminEmail       string
	SuperAdminDisplayName string
	SuperAdminPassword    string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	smtpPort, err := envInt("SMTP_PORT", 587)
	if err != nil {
		return nil, err
	}

	smtpTimeoutSeconds, err := envInt("SMTP_TIMEOUT_SECONDS", 10)
	if err != nil {
		return nil, err
	}

	authTokenTTLHours, err := envInt("AUTH_TOKEN_TTL_HOURS", 12)
	if err != nil {
		return nil, err
	}

	authPasswordResetTokenTTLMinutes, err := envInt("AUTH_PASSWORD_RESET_TOKEN_TTL_MINUTES", 30)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppEnv: os.Getenv("APP_ENV"),
		Port:   os.Getenv("APP_PORT"),
		SMTP: SMTPConfig{
			Host:               envString("SMTP_HOST", "smtp.gmail.com"),
			Port:               smtpPort,
			Username:           strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
			Password:           os.Getenv("SMTP_PASSWORD"),
			AccessToken:        strings.TrimSpace(os.Getenv("SMTP_ACCESS_TOKEN")),
			AuthMethod:         envString("SMTP_AUTH_METHOD", "plain"),
			FromName:           strings.TrimSpace(os.Getenv("SMTP_FROM_NAME")),
			FromAddress:        strings.TrimSpace(os.Getenv("SMTP_FROM_ADDRESS")),
			UseStartTLS:        envBool("SMTP_USE_STARTTLS", true),
			UseImplicitTLS:     envBool("SMTP_USE_IMPLICIT_TLS", false),
			RequireTLS:         envBool("SMTP_REQUIRE_TLS", true),
			InsecureSkipVerify: envBool("SMTP_INSECURE_SKIP_VERIFY", false),
			Timeout:            time.Duration(smtpTimeoutSeconds) * time.Second,
		},
		Auth: AuthConfig{
			JWTSecret:             envString("AUTH_JWT_SECRET", "dev-change-me"),
			TokenTTL:              time.Duration(authTokenTTLHours) * time.Hour,
			PasswordResetTokenTTL: time.Duration(authPasswordResetTokenTTLMinutes) * time.Minute,
			PasswordResetURL:      envString("AUTH_PASSWORD_RESET_URL", "http://localhost:5173/reset-password"),
			CookieName:            envString("AUTH_COOKIE_NAME", "sicou_session"),
			CookiePath:            envString("AUTH_COOKIE_PATH", "/"),
			CookieDomain:          strings.TrimSpace(os.Getenv("AUTH_COOKIE_DOMAIN")),
			CookieSecure:          envBool("AUTH_COOKIE_SECURE", false),
			CookieSameSite:        envString("AUTH_COOKIE_SAMESITE", "Lax"),
			SuperAdminEmail:       envString("SUPERADMIN_EMAIL", "superadmin@sicou.local"),
			SuperAdminDisplayName: envString("SUPERADMIN_DISPLAY_NAME", "Super Administrador"),
			SuperAdminPassword:    envString("SUPERADMIN_PASSWORD", "SuperAdmin123*"),
		},
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.SMTP.FromAddress == "" {
		cfg.SMTP.FromAddress = cfg.SMTP.Username
	}

	return cfg, nil
}

func envString(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}

func envBool(key string, defaultValue bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func envInt(key string, defaultValue int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer: %w", key, err)
	}

	return parsed, nil
}
