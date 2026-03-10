package service

import (
	"context"
	"errors"
	"fmt"
	"html"
	stdmail "net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
)

func (s *AppUserService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	parsedEmail, err := normalizeEmail(req.Email)
	if err != nil {
		return ErrInvalidInput
	}

	if strings.TrimSpace(s.jwtSecret) == "" {
		return ErrTokenNotConfigured
	}

	user, err := s.repo.GetByEmail(ctx, parsedEmail)
	if err != nil {
		return err
	}

	// Avoid email enumeration: the response is the same whether the user exists or not.
	if user == nil || !user.IsActive {
		return nil
	}

	token, expiresAt, err := buildPasswordResetToken(*user, s.jwtSecret, s.passwordResetTokenTTL)
	if err != nil {
		if errors.Is(err, errMissingTokenSecret) {
			return ErrTokenNotConfigured
		}
		return err
	}

	resetLink := buildPasswordResetLink(s.passwordResetFrontendURL, token)
	subject := "Recuperacion de contraseña - SICOU"
	body := buildPasswordResetEmailBody(user.DisplayName, resetLink, expiresAt)

	return s.repo.EnqueuePasswordResetEmail(ctx, user.Email, subject, body)
}

func normalizeEmail(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrInvalidInput
	}

	parsed, err := stdmail.ParseAddress(trimmed)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(parsed.Address), nil
}

func buildPasswordResetLink(baseURL string, token string) string {
	normalizedBaseURL := strings.TrimSpace(baseURL)
	normalizedToken := strings.TrimSpace(token)

	if normalizedToken == "" {
		return normalizedBaseURL
	}
	if normalizedBaseURL == "" {
		return normalizedToken
	}

	parsedURL, err := url.Parse(normalizedBaseURL)
	if err != nil {
		separator := "?"
		if strings.Contains(normalizedBaseURL, "?") {
			separator = "&"
		}
		return normalizedBaseURL + separator + "token=" + url.QueryEscape(normalizedToken)
	}

	query := parsedURL.Query()
	query.Set("token", normalizedToken)
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()
}

func buildPasswordResetEmailBody(displayName string, resetLink string, expiresAt time.Time) string {
	name := strings.TrimSpace(displayName)
	if name == "" {
		name = "usuario"
	}

	escapedName := html.EscapeString(name)
	escapedLink := html.EscapeString(strings.TrimSpace(resetLink))
	expiresLabel := expiresAt.UTC().Format("2006-01-02 15:04:05")

	return fmt.Sprintf(
		"<html><body style=\"font-family: Arial, sans-serif; color: #111827;\">"+
			"<p>Hola %s,</p>"+
			"<p>Recibimos una solicitud para recuperar tu contraseña en SICOU.</p>"+
			"<p><a href=\"%s\" style=\"display:inline-block;padding:10px 16px;background:#1f2937;color:#ffffff;text-decoration:none;border-radius:6px;\">Cambiar contraseña</a></p>"+
			"<p>Este enlace vence el <strong>%s (UTC)</strong>.</p>"+
			"<p>Si no solicitaste este cambio, ignora este mensaje.</p>"+
			"<p style=\"font-size:12px;color:#6b7280;\">Si el boton no funciona, copia y pega este enlace:<br>%s</p>"+
			"</body></html>",
		escapedName,
		escapedLink,
		expiresLabel,
		escapedLink,
	)
}
