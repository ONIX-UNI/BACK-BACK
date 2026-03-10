package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type sessionState struct {
	user      *dto.AppUser
	roles     []string
	expiresAt time.Time
	signature string
}

func (s *AppUserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, string, error) {
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)
	if email == "" || password == "" {
		return nil, "", ErrInvalidInput
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, "", ErrUserInactive
	}
	if !isPasswordValid(user.PasswordHash, password) {
		return nil, "", ErrInvalidCredentials
	}

	lastAccessAt := time.Now().UTC()
	if err := s.repo.UpdateLastAccess(ctx, user.ID, lastAccessAt); err != nil {
		return nil, "", err
	}
	user.LastAccessAt = &lastAccessAt

	roles, err := s.repo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, "", err
	}

	token, expiresAt, err := buildAccessToken(*user, roles, s.jwtSecret, s.tokenTTL)
	if err != nil {
		if errors.Is(err, errMissingTokenSecret) {
			return nil, "", ErrTokenNotConfigured
		}
		return nil, "", err
	}

	return buildSessionResponse(*user, roles, expiresAt), token, nil
}

func (s *AppUserService) Me(ctx context.Context, rawToken string) (*dto.LoginResponse, error) {
	sessionState, err := s.resolveSession(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	return buildSessionResponse(*sessionState.user, sessionState.roles, sessionState.expiresAt), nil
}

func (s *AppUserService) Refresh(ctx context.Context, rawToken string) (*dto.LoginResponse, string, error) {
	sessionState, err := s.resolveSession(ctx, rawToken)
	if err != nil {
		return nil, "", err
	}

	lastAccessAt := time.Now().UTC()
	if err := s.repo.UpdateLastAccess(ctx, sessionState.user.ID, lastAccessAt); err != nil {
		return nil, "", err
	}
	sessionState.user.LastAccessAt = &lastAccessAt

	if err := s.repo.RevokeToken(ctx, sessionState.signature, sessionState.expiresAt); err != nil {
		return nil, "", err
	}

	newToken, newExpiresAt, err := buildAccessToken(*sessionState.user, sessionState.roles, s.jwtSecret, s.tokenTTL)
	if err != nil {
		if errors.Is(err, errMissingTokenSecret) {
			return nil, "", ErrTokenNotConfigured
		}
		return nil, "", err
	}

	return buildSessionResponse(*sessionState.user, sessionState.roles, newExpiresAt), newToken, nil
}

func (s *AppUserService) Logout(ctx context.Context, rawToken string) error {
	if strings.TrimSpace(rawToken) == "" {
		return nil
	}

	_, signature, expiresAt, err := parseAccessTokenAllowExpired(rawToken, s.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, errMissingTokenSecret):
			return ErrTokenNotConfigured
		case errors.Is(err, errInvalidToken), errors.Is(err, errExpiredToken):
			return nil
		default:
			return err
		}
	}

	if expiresAt.After(time.Now().UTC()) {
		if err := s.repo.RevokeToken(ctx, signature, expiresAt); err != nil {
			return err
		}
	}

	return nil
}

func (s *AppUserService) resolveSession(ctx context.Context, rawToken string) (*sessionState, error) {
	claims, signature, expiresAt, err := parseAccessToken(rawToken, s.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, errMissingTokenSecret):
			return nil, ErrTokenNotConfigured
		case errors.Is(err, errExpiredToken):
			return nil, ErrSessionExpired
		case errors.Is(err, errInvalidToken):
			return nil, ErrSessionInvalid
		default:
			return nil, err
		}
	}

	revoked, err := s.repo.IsTokenRevoked(ctx, signature)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, ErrSessionRevoked
	}

	userID, err := uuid.Parse(strings.TrimSpace(claims.Subject))
	if err != nil {
		return nil, ErrSessionInvalid
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrSessionInvalid
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	roles, err := s.repo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &sessionState{
		user:      user,
		roles:     roles,
		expiresAt: expiresAt,
		signature: signature,
	}, nil
}

func buildSessionResponse(user dto.AppUser, roles []string, expiresAt time.Time) *dto.LoginResponse {
	expiresIn := int64(time.Until(expiresAt).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	return &dto.LoginResponse{
		ExpiresAt: expiresAt,
		ExpiresIn: expiresIn,
		User: dto.LoginUserResponse{
			ID:          user.ID,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			Roles:       roles,
		},
	}
}

func isPasswordValid(storedHash, rawPassword string) bool {
	if strings.TrimSpace(storedHash) == "" || strings.TrimSpace(rawPassword) == "" {
		return false
	}

	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(rawPassword)) == nil {
		return true
	}

	return storedHash == rawPassword
}
