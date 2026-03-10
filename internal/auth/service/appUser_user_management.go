package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/DuvanRozoParra/sicou/internal/auth/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AppUserService) Create(ctx context.Context, req dto.CreateAppUserRequest, actorUserID *uuid.UUID) (*dto.AppUser, error) {
	req.Email = strings.TrimSpace(req.Email)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.PasswordHash = strings.TrimSpace(req.PasswordHash)
	req.Role = dto.NormalizeRoleCode(req.Role)

	if req.Email == "" ||
		req.DisplayName == "" ||
		req.PasswordHash == "" ||
		req.Role == "" {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyUsed
	}

	if !looksLikeBcryptHash(req.PasswordHash) {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return nil, ErrPasswordHashFailure
		}
		req.PasswordHash = string(passwordHash)
	}

	user, err := s.repo.Create(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrRoleNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	if err := s.attachRoles(ctx, user); err != nil {
		return nil, err
	}

	if err := s.createUserAuditLog(ctx, dto.UserAuditLogEntry{
		ActorUserID:     actorUserID,
		Action:          dto.UserAuditActionCreated,
		TargetUserID:    pointerUUID(user.ID),
		TargetUserName:  strings.TrimSpace(user.DisplayName),
		TargetUserEmail: strings.TrimSpace(user.Email),
		Detail:          buildCreateUserDetail(user.DisplayName, user.Email, user.Role),
		Metadata: dto.UserAuditMetadata{
			After: userAuditState(*user),
		},
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AppUserService) GetByID(ctx context.Context, id uuid.UUID) (*dto.AppUser, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if err := s.attachRoles(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AppUserService) GetByEmail(ctx context.Context, email string) (*dto.AppUser, error) {
	if strings.TrimSpace(email) == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if err := s.attachRoles(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AppUserService) EnsureAuthStorage(ctx context.Context) error {
	return s.repo.EnsureAuthStorage(ctx)
}

func (s *AppUserService) GetByDisplayName(ctx context.Context, displayName string, limit, offset int) (*dto.ListAppUsersResponse, error) {
	if strings.TrimSpace(displayName) == "" {
		return nil, ErrInvalidInput
	}

	result, err := s.repo.GetByDisplayName(ctx, displayName, limit, offset)
	if err != nil {
		return nil, err
	}
	for i := range result.Items {
		if err := s.attachRoles(ctx, &result.Items[i]); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *AppUserService) List(ctx context.Context, limit, offset int) (*dto.ListAppUsersResponse, error) {
	result, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	for i := range result.Items {
		if err := s.attachRoles(ctx, &result.Items[i]); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *AppUserService) ListWithFilters(ctx context.Context, req dto.ListAppUsersFilterRequest) (*dto.ListAppUsersResponse, error) {
	role := dto.NormalizeRoleCode(req.Role)
	if role != "" {
		exists, err := s.repo.RoleExists(ctx, role)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrInvalidInput
		}
	}

	result, err := s.repo.ListWithFilters(ctx, dto.ListAppUsersFilterRequest{
		Query:  strings.TrimSpace(req.Query),
		Role:   role,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}
	for i := range result.Items {
		if err := s.attachRoles(ctx, &result.Items[i]); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *AppUserService) ListOptions(ctx context.Context, req dto.ListAppUserOptionsRequest) (*dto.ListAppUserOptionsResponse, error) {
	normalizedRoles := dto.NormalizeRoleCodes(req.Roles)
	if len(normalizedRoles) == 0 {
		return nil, ErrInvalidInput
	}

	result, err := s.repo.ListOptions(ctx, dto.ListAppUserOptionsRequest{
		Roles:    normalizedRoles,
		Search:   strings.TrimSpace(req.Search),
		IsActive: req.IsActive,
		Limit:    req.Limit,
		Offset:   req.Offset,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AppUserService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateAppUserRequest, actorUserID *uuid.UUID) (*dto.AppUser, error) {
	if req.DisplayName == nil && req.IsActive == nil && req.Role == nil {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrUserNotFound
	}
	if err := s.attachRoles(ctx, existing); err != nil {
		return nil, err
	}

	beforeState := userAuditState(*existing)
	beforeRole := strings.TrimSpace(existing.Role)
	beforeIsActive := existing.IsActive

	if req.Role != nil {
		role := dto.NormalizeRoleCode(*req.Role)
		if role == "" {
			return nil, ErrInvalidInput
		}
		if err := s.repo.ReplacePrimaryRole(ctx, id, role); err != nil {
			if errors.Is(err, models.ErrRoleNotFound) {
				return nil, ErrRoleNotFound
			}
			return nil, err
		}
	}

	var updated *dto.AppUser
	if req.DisplayName != nil || req.IsActive != nil {
		updated, err = s.repo.Update(ctx, id, req)
		if err != nil {
			return nil, err
		}
	} else {
		updated, err = s.repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if updated == nil {
			return nil, ErrUserNotFound
		}
	}
	if err := s.attachRoles(ctx, updated); err != nil {
		return nil, err
	}

	afterState := userAuditState(*updated)
	afterRole := strings.TrimSpace(updated.Role)
	afterIsActive := updated.IsActive

	events := buildUserUpdateAuditEvents(
		*existing,
		*updated,
		beforeState,
		afterState,
		beforeRole,
		afterRole,
		beforeIsActive,
		afterIsActive,
	)
	for _, event := range events {
		event.ActorUserID = actorUserID
		event.TargetUserID = pointerUUID(updated.ID)
		event.TargetUserName = strings.TrimSpace(updated.DisplayName)
		event.TargetUserEmail = strings.TrimSpace(updated.Email)
		if err := s.createUserAuditLog(ctx, event); err != nil {
			return nil, err
		}
	}

	return updated, nil
}

func (s *AppUserService) EnsureSingleSuperAdmin(ctx context.Context, req dto.BootstrapSuperAdminRequest) error {
	email := strings.TrimSpace(req.Email)
	displayName := strings.TrimSpace(req.DisplayName)
	password := strings.TrimSpace(req.Password)
	if email == "" || displayName == "" || password == "" {
		return ErrInvalidInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrPasswordHashFailure
	}

	return s.repo.EnsureSingleSuperAdmin(ctx, email, displayName, string(passwordHash))
}

func (s *AppUserService) Delete(ctx context.Context, id uuid.UUID, actorUserID *uuid.UUID) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}
	if err := s.attachRoles(ctx, existing); err != nil {
		return err
	}

	beforeState := userAuditState(*existing)

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return s.createUserAuditLog(ctx, dto.UserAuditLogEntry{
		ActorUserID:     actorUserID,
		Action:          dto.UserAuditActionDeleted,
		TargetUserID:    pointerUUID(existing.ID),
		TargetUserName:  strings.TrimSpace(existing.DisplayName),
		TargetUserEmail: strings.TrimSpace(existing.Email),
		Detail:          buildDeleteUserDetail(existing.DisplayName, existing.Email),
		Metadata: dto.UserAuditMetadata{
			Before: beforeState,
		},
	})
}

func (s *AppUserService) attachRoles(ctx context.Context, user *dto.AppUser) error {
	roles, err := s.repo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	user.Roles = roles
	user.Role = ""
	if len(roles) > 0 {
		user.Role = roles[0]
	}

	return nil
}

func looksLikeBcryptHash(value string) bool {
	trimmed := strings.TrimSpace(value)
	return strings.HasPrefix(trimmed, "$2a$") ||
		strings.HasPrefix(trimmed, "$2b$") ||
		strings.HasPrefix(trimmed, "$2y$")
}
