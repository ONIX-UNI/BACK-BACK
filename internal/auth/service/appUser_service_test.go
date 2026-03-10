package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/DuvanRozoParra/sicou/internal/auth/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockAppUserRepository struct {
	createFn                 func(ctx context.Context, req dto.CreateAppUserRequest) (*dto.AppUser, error)
	getByIDFn                func(ctx context.Context, id uuid.UUID) (*dto.AppUser, error)
	getByEmailFn             func(ctx context.Context, email string) (*dto.AppUser, error)
	enqueuePasswordResetEmailFn func(ctx context.Context, toEmail string, subject string, body string) error
	updatePasswordFn         func(ctx context.Context, id uuid.UUID, passwordHash string) error
	getRolesByUserIDFn       func(ctx context.Context, userID uuid.UUID) ([]string, error)
	roleExistsFn             func(ctx context.Context, roleCode string) (bool, error)
	getByDisplayNameFn       func(ctx context.Context, displayName string, limit, offset int) (*dto.ListAppUsersResponse, error)
	listWithFiltersFn        func(ctx context.Context, req dto.ListAppUsersFilterRequest) (*dto.ListAppUsersResponse, error)
	listOptionsFn            func(ctx context.Context, req dto.ListAppUserOptionsRequest) (*dto.ListAppUserOptionsResponse, error)
	listFn                   func(ctx context.Context, limit, offset int) (*dto.ListAppUsersResponse, error)
	updateFn                 func(ctx context.Context, id uuid.UUID, req dto.UpdateAppUserRequest) (*dto.AppUser, error)
	replacePrimaryRoleFn     func(ctx context.Context, userID uuid.UUID, roleCode string) error
	ensureSingleSuperAdminFn func(ctx context.Context, email, displayName, passwordHash string) error
	ensureAuthStorageFn      func(ctx context.Context) error
	isTokenRevokedFn         func(ctx context.Context, tokenSignature string) (bool, error)
	revokeTokenFn            func(ctx context.Context, tokenSignature string, expiresAt time.Time) error
	createUserAuditLogFn     func(ctx context.Context, entry dto.UserAuditLogEntry) error
	listUserAuditLogFn       func(ctx context.Context, req dto.ListUserAuditLogRequest) (*dto.ListUserAuditLogResponse, error)
	updateLastAccessFn       func(ctx context.Context, id uuid.UUID, accessedAt time.Time) error
	deleteFn                 func(ctx context.Context, id uuid.UUID) error
}

func (m mockAppUserRepository) Create(ctx context.Context, req dto.CreateAppUserRequest) (*dto.AppUser, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil, nil
}

func (m mockAppUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*dto.AppUser, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m mockAppUserRepository) GetByEmail(ctx context.Context, email string) (*dto.AppUser, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m mockAppUserRepository) EnqueuePasswordResetEmail(ctx context.Context, toEmail string, subject string, body string) error {
	if m.enqueuePasswordResetEmailFn != nil {
		return m.enqueuePasswordResetEmailFn(ctx, toEmail, subject, body)
	}
	return nil
}

func (m mockAppUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	if m.updatePasswordFn != nil {
		return m.updatePasswordFn(ctx, id, passwordHash)
	}
	return nil
}

func (m mockAppUserRepository) GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	if m.getRolesByUserIDFn != nil {
		return m.getRolesByUserIDFn(ctx, userID)
	}
	return nil, nil
}

func (m mockAppUserRepository) RoleExists(ctx context.Context, roleCode string) (bool, error) {
	if m.roleExistsFn != nil {
		return m.roleExistsFn(ctx, roleCode)
	}
	return false, nil
}

func (m mockAppUserRepository) GetByDisplayName(ctx context.Context, displayName string, limit, offset int) (*dto.ListAppUsersResponse, error) {
	if m.getByDisplayNameFn != nil {
		return m.getByDisplayNameFn(ctx, displayName, limit, offset)
	}
	return &dto.ListAppUsersResponse{}, nil
}

func (m mockAppUserRepository) ListWithFilters(ctx context.Context, req dto.ListAppUsersFilterRequest) (*dto.ListAppUsersResponse, error) {
	if m.listWithFiltersFn != nil {
		return m.listWithFiltersFn(ctx, req)
	}
	return &dto.ListAppUsersResponse{}, nil
}

func (m mockAppUserRepository) List(ctx context.Context, limit, offset int) (*dto.ListAppUsersResponse, error) {
	if m.listFn != nil {
		return m.listFn(ctx, limit, offset)
	}
	return &dto.ListAppUsersResponse{}, nil
}

func (m mockAppUserRepository) ListOptions(ctx context.Context, req dto.ListAppUserOptionsRequest) (*dto.ListAppUserOptionsResponse, error) {
	if m.listOptionsFn != nil {
		return m.listOptionsFn(ctx, req)
	}
	return &dto.ListAppUserOptionsResponse{}, nil
}

func (m mockAppUserRepository) Update(ctx context.Context, id uuid.UUID, req dto.UpdateAppUserRequest) (*dto.AppUser, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, req)
	}
	return nil, nil
}

func (m mockAppUserRepository) ReplacePrimaryRole(ctx context.Context, userID uuid.UUID, roleCode string) error {
	if m.replacePrimaryRoleFn != nil {
		return m.replacePrimaryRoleFn(ctx, userID, roleCode)
	}
	return nil
}

func (m mockAppUserRepository) EnsureSingleSuperAdmin(ctx context.Context, email, displayName, passwordHash string) error {
	if m.ensureSingleSuperAdminFn != nil {
		return m.ensureSingleSuperAdminFn(ctx, email, displayName, passwordHash)
	}
	return nil
}

func (m mockAppUserRepository) EnsureAuthStorage(ctx context.Context) error {
	if m.ensureAuthStorageFn != nil {
		return m.ensureAuthStorageFn(ctx)
	}
	return nil
}

func (m mockAppUserRepository) IsTokenRevoked(ctx context.Context, tokenSignature string) (bool, error) {
	if m.isTokenRevokedFn != nil {
		return m.isTokenRevokedFn(ctx, tokenSignature)
	}
	return false, nil
}

func (m mockAppUserRepository) RevokeToken(ctx context.Context, tokenSignature string, expiresAt time.Time) error {
	if m.revokeTokenFn != nil {
		return m.revokeTokenFn(ctx, tokenSignature, expiresAt)
	}
	return nil
}

func (m mockAppUserRepository) CreateUserAuditLog(ctx context.Context, entry dto.UserAuditLogEntry) error {
	if m.createUserAuditLogFn != nil {
		return m.createUserAuditLogFn(ctx, entry)
	}
	return nil
}

func (m mockAppUserRepository) ListUserAuditLog(ctx context.Context, req dto.ListUserAuditLogRequest) (*dto.ListUserAuditLogResponse, error) {
	if m.listUserAuditLogFn != nil {
		return m.listUserAuditLogFn(ctx, req)
	}
	return &dto.ListUserAuditLogResponse{}, nil
}

func (m mockAppUserRepository) UpdateLastAccess(ctx context.Context, id uuid.UUID, accessedAt time.Time) error {
	if m.updateLastAccessFn != nil {
		return m.updateLastAccessFn(ctx, id, accessedAt)
	}
	return nil
}

func (m mockAppUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func TestLoginSuccess(t *testing.T) {
	password := "Admin123*"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to build hash: %v", err)
	}

	userID := uuid.New()
	svc := NewAppUserService(mockAppUserRepository{
		getByEmailFn: func(ctx context.Context, email string) (*dto.AppUser, error) {
			return &dto.AppUser{
				ID:           userID,
				Email:        "admin@example.com",
				DisplayName:  "Admin",
				PasswordHash: string(passwordHash),
				IsActive:     true,
			}, nil
		},
		getRolesByUserIDFn: func(ctx context.Context, userID uuid.UUID) ([]string, error) {
			return []string{"SUPER_ADMIN"}, nil
		},
	}, "test-secret", time.Hour)

	response, token, err := svc.Login(context.Background(), dto.LoginRequest{
		Email:    "admin@example.com",
		Password: password,
	})
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty session token")
	}
	if response.ExpiresIn <= 0 {
		t.Fatalf("expected positive expires_in, got %d", response.ExpiresIn)
	}
	if response.User.ID != userID {
		t.Fatalf("expected user id %s, got %s", userID, response.User.ID)
	}
	if len(response.User.Roles) != 1 || response.User.Roles[0] != "SUPER_ADMIN" {
		t.Fatalf("expected SUPER_ADMIN role, got %#v", response.User.Roles)
	}
}

func TestLoginUpdatesLastAccess(t *testing.T) {
	password := "Admin123*"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to build hash: %v", err)
	}

	userID := uuid.New()
	updatedLastAccess := false
	svc := NewAppUserService(mockAppUserRepository{
		getByEmailFn: func(ctx context.Context, email string) (*dto.AppUser, error) {
			return &dto.AppUser{
				ID:           userID,
				Email:        "admin@example.com",
				DisplayName:  "Admin",
				PasswordHash: string(passwordHash),
				IsActive:     true,
			}, nil
		},
		getRolesByUserIDFn: func(ctx context.Context, userID uuid.UUID) ([]string, error) {
			return []string{"SUPER_ADMIN"}, nil
		},
		updateLastAccessFn: func(ctx context.Context, id uuid.UUID, accessedAt time.Time) error {
			if id != userID {
				t.Fatalf("expected user id %s, got %s", userID, id)
			}
			if accessedAt.IsZero() {
				t.Fatal("expected non-zero accessedAt")
			}
			updatedLastAccess = true
			return nil
		},
	}, "test-secret", time.Hour)

	_, _, err = svc.Login(context.Background(), dto.LoginRequest{
		Email:    "admin@example.com",
		Password: password,
	})
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}
	if !updatedLastAccess {
		t.Fatal("expected UpdateLastAccess to be called")
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to build hash: %v", err)
	}

	svc := NewAppUserService(mockAppUserRepository{
		getByEmailFn: func(ctx context.Context, email string) (*dto.AppUser, error) {
			return &dto.AppUser{
				ID:           uuid.New(),
				Email:        "admin@example.com",
				DisplayName:  "Admin",
				PasswordHash: string(passwordHash),
				IsActive:     true,
			}, nil
		},
		getRolesByUserIDFn: func(ctx context.Context, userID uuid.UUID) ([]string, error) {
			return []string{"SUPER_ADMIN"}, nil
		},
	}, "test-secret", time.Hour)

	_, _, err = svc.Login(context.Background(), dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestCreateRequiresRole(t *testing.T) {
	svc := NewAppUserService(mockAppUserRepository{}, "test-secret", time.Hour)

	_, err := svc.Create(context.Background(), dto.CreateAppUserRequest{
		Email:        "admin@example.com",
		DisplayName:  "Admin",
		PasswordHash: "secret",
	}, nil)
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreateMapsRoleNotFound(t *testing.T) {
	svc := NewAppUserService(mockAppUserRepository{
		getByEmailFn: func(ctx context.Context, email string) (*dto.AppUser, error) {
			return nil, nil
		},
		createFn: func(ctx context.Context, req dto.CreateAppUserRequest) (*dto.AppUser, error) {
			return nil, models.ErrRoleNotFound
		},
	}, "test-secret", time.Hour)

	_, err := svc.Create(context.Background(), dto.CreateAppUserRequest{
		Email:        "admin@example.com",
		DisplayName:  "Admin",
		PasswordHash: "secret",
		Role:         "UNKNOWN",
	}, nil)
	if !errors.Is(err, ErrRoleNotFound) {
		t.Fatalf("expected ErrRoleNotFound, got %v", err)
	}
}

func TestCreateReturnsAssignedRoles(t *testing.T) {
	userID := uuid.New()
	var capturedRole string
	var capturedPasswordHash string

	svc := NewAppUserService(mockAppUserRepository{
		getByEmailFn: func(ctx context.Context, email string) (*dto.AppUser, error) {
			return nil, nil
		},
		createFn: func(ctx context.Context, req dto.CreateAppUserRequest) (*dto.AppUser, error) {
			capturedRole = req.Role
			capturedPasswordHash = req.PasswordHash
			return &dto.AppUser{
				ID:          userID,
				Email:       req.Email,
				DisplayName: req.DisplayName,
				IsActive:    true,
			}, nil
		},
		getRolesByUserIDFn: func(ctx context.Context, userID uuid.UUID) ([]string, error) {
			return []string{"SECRETARIA"}, nil
		},
	}, "test-secret", time.Hour)

	user, err := svc.Create(context.Background(), dto.CreateAppUserRequest{
		Email:        "admin@example.com",
		DisplayName:  "Admin",
		PasswordHash: "secret",
		Role:         "secretaria",
	}, nil)
	if err != nil {
		t.Fatalf("expected create success, got %v", err)
	}
	if capturedRole != "SECRETARIA" {
		t.Fatalf("expected normalized role SECRETARIA, got %q", capturedRole)
	}
	if capturedPasswordHash == "secret" {
		t.Fatal("expected password to be hashed before repository create")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(capturedPasswordHash), []byte("secret")); err != nil {
		t.Fatalf("expected bcrypt hash for password, got error: %v", err)
	}
	if len(user.Roles) != 1 || user.Roles[0] != "SECRETARIA" {
		t.Fatalf("expected assigned role SECRETARIA, got %#v", user.Roles)
	}
	if user.Role != "SECRETARIA" {
		t.Fatalf("expected role SECRETARIA, got %q", user.Role)
	}
}

func TestEnsureSingleSuperAdminHashesPassword(t *testing.T) {
	var capturedHash string
	svc := NewAppUserService(mockAppUserRepository{
		ensureSingleSuperAdminFn: func(ctx context.Context, email, displayName, passwordHash string) error {
			if email != "admin@example.com" {
				t.Fatalf("unexpected email: %q", email)
			}
			if displayName != "Super Admin" {
				t.Fatalf("unexpected displayName: %q", displayName)
			}
			capturedHash = passwordHash
			return nil
		},
	}, "test-secret", time.Hour)

	err := svc.EnsureSingleSuperAdmin(context.Background(), dto.BootstrapSuperAdminRequest{
		Email:       "admin@example.com",
		DisplayName: "Super Admin",
		Password:    "Admin123*",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedHash == "" {
		t.Fatal("expected password hash to be captured")
	}
	if capturedHash == "Admin123*" {
		t.Fatal("expected bcrypt hash, got plain password")
	}
	if bcrypt.CompareHashAndPassword([]byte(capturedHash), []byte("Admin123*")) != nil {
		t.Fatal("captured hash does not match the original password")
	}
}

func TestListOptionsRequiresRoles(t *testing.T) {
	svc := NewAppUserService(mockAppUserRepository{}, "test-secret", time.Hour)

	_, err := svc.ListOptions(context.Background(), dto.ListAppUserOptionsRequest{
		IsActive: true,
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestListOptionsNormalizesAndPaginates(t *testing.T) {
	var capturedReq dto.ListAppUserOptionsRequest

	itemID := uuid.New()
	svc := NewAppUserService(mockAppUserRepository{
		listOptionsFn: func(ctx context.Context, req dto.ListAppUserOptionsRequest) (*dto.ListAppUserOptionsResponse, error) {
			capturedReq = req
			return &dto.ListAppUserOptionsResponse{
				Items: []dto.AppUserOptionItem{
					{
						ID:          itemID,
						DisplayName: "Juan Perez",
						Email:       "juan@unimeta.edu.co",
						Role:        "ESTUDIANTE",
					},
				},
				Total:   248,
				Limit:   req.Limit,
				Offset:  req.Offset,
				HasNext: true,
			}, nil
		},
	}, "test-secret", time.Hour)

	result, err := svc.ListOptions(context.Background(), dto.ListAppUserOptionsRequest{
		Roles:    []string{" estudiante ", "COORDINADOR", "ESTUDIANTE"},
		Search:   " juan ",
		IsActive: true,
		Limit:    0,
		Offset:   -5,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capturedReq.Limit != 0 {
		t.Fatalf("expected raw limit passthrough 0, got %d", capturedReq.Limit)
	}
	if capturedReq.Offset != -5 {
		t.Fatalf("expected raw offset passthrough -5, got %d", capturedReq.Offset)
	}
	if capturedReq.Search != "juan" {
		t.Fatalf("expected trimmed search 'juan', got %q", capturedReq.Search)
	}
	if len(capturedReq.Roles) != 2 || capturedReq.Roles[0] != "ESTUDIANTE" || capturedReq.Roles[1] != "COORDINADOR" {
		t.Fatalf("expected normalized roles [ESTUDIANTE COORDINADOR], got %#v", capturedReq.Roles)
	}
	if result.Total != 248 {
		t.Fatalf("expected total 248, got %d", result.Total)
	}
	if !result.HasNext {
		t.Fatal("expected has_next to be true")
	}
	if len(result.Items) != 1 || result.Items[0].ID != itemID {
		t.Fatalf("unexpected items: %#v", result.Items)
	}
}

func TestListWithFiltersRejectsUnknownRole(t *testing.T) {
	svc := NewAppUserService(mockAppUserRepository{
		roleExistsFn: func(ctx context.Context, roleCode string) (bool, error) {
			return false, nil
		},
	}, "test-secret", time.Hour)

	_, err := svc.ListWithFilters(context.Background(), dto.ListAppUsersFilterRequest{
		Query:  "dani",
		Role:   "SECRETARIA",
		Limit:  10,
		Offset: 0,
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestListWithFiltersNormalizesAndBuildsHasNext(t *testing.T) {
	userID := uuid.New()
	var capturedReq dto.ListAppUsersFilterRequest

	svc := NewAppUserService(mockAppUserRepository{
		roleExistsFn: func(ctx context.Context, roleCode string) (bool, error) {
			return roleCode == "SECRETARIA", nil
		},
		listWithFiltersFn: func(ctx context.Context, req dto.ListAppUsersFilterRequest) (*dto.ListAppUsersResponse, error) {
			capturedReq = req
			return &dto.ListAppUsersResponse{
				Items: []dto.AppUser{
					{
						ID:          userID,
						DisplayName: "Daniel Perez",
						Email:       "daniel@x.com",
						IsActive:    true,
					},
				},
				Total:       37,
				ActiveCount: 31,
				Limit:       req.Limit,
				Offset:      req.Offset,
				HasNext:     true,
			}, nil
		},
		getRolesByUserIDFn: func(ctx context.Context, userID uuid.UUID) ([]string, error) {
			return []string{"SECRETARIA"}, nil
		},
	}, "test-secret", time.Hour)

	result, err := svc.ListWithFilters(context.Background(), dto.ListAppUsersFilterRequest{
		Query:  " dani ",
		Role:   " secretaria ",
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capturedReq.Query != "dani" {
		t.Fatalf("expected trimmed query 'dani', got %q", capturedReq.Query)
	}
	if capturedReq.Role != "SECRETARIA" {
		t.Fatalf("expected normalized role SECRETARIA, got %q", capturedReq.Role)
	}
	if result.Total != 37 || result.ActiveCount != 31 {
		t.Fatalf("unexpected totals: total=%d active_count=%d", result.Total, result.ActiveCount)
	}
	if !result.HasNext {
		t.Fatal("expected has_next to be true")
	}
	if len(result.Items) != 1 || result.Items[0].Role != "SECRETARIA" {
		t.Fatalf("expected one item with role SECRETARIA, got %#v", result.Items)
	}
}
