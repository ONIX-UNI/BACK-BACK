package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AppUserService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	rawToken := strings.TrimSpace(req.Token)
	rawPassword := strings.TrimSpace(req.Password)
	if rawToken == "" || rawPassword == "" {
		return ErrInvalidInput
	}

	claims, signature, expiresAt, err := parsePasswordResetToken(rawToken, s.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, errMissingTokenSecret):
			return ErrTokenNotConfigured
		case errors.Is(err, errExpiredToken):
			return ErrPasswordResetTokenExpired
		case errors.Is(err, errInvalidToken):
			return ErrPasswordResetTokenInvalid
		default:
			return err
		}
	}

	revoked, err := s.repo.IsTokenRevoked(ctx, signature)
	if err != nil {
		return err
	}
	if revoked {
		return ErrPasswordResetTokenInvalid
	}

	userID, err := uuid.Parse(strings.TrimSpace(claims.Subject))
	if err != nil {
		return ErrPasswordResetTokenInvalid
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrPasswordResetTokenInvalid
	}
	if !user.IsActive {
		return ErrUserInactive
	}
	if !strings.EqualFold(strings.TrimSpace(user.Email), strings.TrimSpace(claims.Email)) {
		return ErrPasswordResetTokenInvalid
	}

	passwordHash := rawPassword
	if !looksLikeBcryptHash(passwordHash) {
		generatedHash, hashErr := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
		if hashErr != nil {
			return ErrPasswordHashFailure
		}
		passwordHash = string(generatedHash)
	}

	if err := s.repo.UpdatePassword(ctx, user.ID, passwordHash); err != nil {
		return err
	}

	if err := s.repo.RevokeToken(ctx, signature, expiresAt); err != nil {
		return err
	}

	return s.createUserAuditLog(ctx, dto.UserAuditLogEntry{
		Action:          dto.UserAuditActionPasswordReset,
		TargetUserID:    pointerUUID(user.ID),
		TargetUserName:  strings.TrimSpace(user.DisplayName),
		TargetUserEmail: strings.TrimSpace(user.Email),
		Detail:          "Restablecio su contrasena mediante recuperacion",
	})
}
