package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/appointments/repository"
	"github.com/google/uuid"
)

type assignmentRepoMock struct {
	createFn            func(ctx context.Context, in repository.CreateRecord) (repository.CreateResult, error)
	listFn              func(ctx context.Context, in repository.ListInput) (repository.ListResult, error)
	listAssignedFn      func(ctx context.Context, in repository.ListAssignedInput) (repository.ListResult, error)
	assignPreturnoFn    func(ctx context.Context, in repository.AssignPreturnoInput) (repository.AssignPreturnoResult, error)
	assignmentOptionsFn func(ctx context.Context) (repository.AssignmentOptionsResult, error)
}

func (m assignmentRepoMock) Create(ctx context.Context, in repository.CreateRecord) (repository.CreateResult, error) {
	if m.createFn != nil {
		return m.createFn(ctx, in)
	}
	return repository.CreateResult{}, nil
}

func (m assignmentRepoMock) List(ctx context.Context, in repository.ListInput) (repository.ListResult, error) {
	if m.listFn != nil {
		return m.listFn(ctx, in)
	}
	return repository.ListResult{}, nil
}

func (m assignmentRepoMock) ListAssigned(ctx context.Context, in repository.ListAssignedInput) (repository.ListResult, error) {
	if m.listAssignedFn != nil {
		return m.listAssignedFn(ctx, in)
	}
	return repository.ListResult{}, nil
}

func (m assignmentRepoMock) AssignPreturno(ctx context.Context, in repository.AssignPreturnoInput) (repository.AssignPreturnoResult, error) {
	if m.assignPreturnoFn != nil {
		return m.assignPreturnoFn(ctx, in)
	}
	return repository.AssignPreturnoResult{}, nil
}

func (m assignmentRepoMock) AssignmentOptions(ctx context.Context) (repository.AssignmentOptionsResult, error) {
	if m.assignmentOptionsFn != nil {
		return m.assignmentOptionsFn(ctx)
	}
	return repository.AssignmentOptionsResult{}, nil
}

func TestAssignPreturno_ValidatesRequiredFields(t *testing.T) {
	svc := NewService(assignmentRepoMock{}, nil)

	_, err := svc.AssignPreturno(context.Background(), AssignPreturnoInput{
		PreturnoID: uuid.NewString(),
		ActorRoles: []string{"COORDINADOR"},
	})

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got: %v", err)
	}
	if got := validationErr.Fields["coordinator_id"]; got == "" {
		t.Fatalf("expected coordinator_id validation error, got: %+v", validationErr.Fields)
	}
	if got := validationErr.Fields["service_type_id"]; got == "" {
		t.Fatalf("expected service_type_id validation error, got: %+v", validationErr.Fields)
	}
}

func TestAssignPreturno_ReturnsForbiddenWhenRoleDoesNotAllowAssignment(t *testing.T) {
	svc := NewService(assignmentRepoMock{}, nil)

	_, err := svc.AssignPreturno(context.Background(), AssignPreturnoInput{
		PreturnoID:    uuid.NewString(),
		CoordinatorID: uuid.NewString(),
		ServiceTypeID: "1",
		ActorRoles:    []string{"ESTUDIANTE"},
	})
	if !errors.Is(err, ErrAssignmentForbidden) {
		t.Fatalf("expected ErrAssignmentForbidden, got: %v", err)
	}
}

func TestAssignPreturno_AssignsAndMapsStatus(t *testing.T) {
	preturnoID := uuid.NewString()
	coordinatorID := uuid.NewString()
	actorID := uuid.NewString()
	now := time.Now().UTC()

	var captured repository.AssignPreturnoInput
	svc := NewService(assignmentRepoMock{
		assignPreturnoFn: func(ctx context.Context, in repository.AssignPreturnoInput) (repository.AssignPreturnoResult, error) {
			captured = in
			return repository.AssignPreturnoResult{
				ID:                    preturnoID,
				Status:                repository.StatusAsignadoPreturno,
				AssignedCoordinatorID: coordinatorID,
				ServiceTypeID:         2,
				TimelineEvent: repository.AssignmentTimelineEvent{
					ID:        uuid.NewString(),
					Title:     "Asignacion de turno",
					Detail:    "Asignacion de turno a Juan Perez con una nota: Llamar hoy",
					CreatedAt: now,
				},
			}, nil
		},
	}, nil)

	result, err := svc.AssignPreturno(context.Background(), AssignPreturnoInput{
		PreturnoID:    preturnoID,
		CoordinatorID: coordinatorID,
		ServiceTypeID: "2",
		Observations:  "Llamar hoy",
		ActorUserID:   actorID,
		ActorRoles:    []string{"COORDINADOR"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "asignado" {
		t.Fatalf("expected status asignado, got %q", result.Status)
	}
	if result.ServiceTypeID != "2" {
		t.Fatalf("expected service_type_id 2, got %q", result.ServiceTypeID)
	}
	if result.AssignedCoordinatorID != coordinatorID {
		t.Fatalf("unexpected assigned coordinator id: %s", result.AssignedCoordinatorID)
	}
	if result.TimelineEvent.CreatedAt != now {
		t.Fatalf("unexpected timeline created_at: %s", result.TimelineEvent.CreatedAt)
	}

	if captured.PreturnoID != preturnoID {
		t.Fatalf("unexpected repository preturno id: %s", captured.PreturnoID)
	}
	if captured.CoordinatorID != coordinatorID {
		t.Fatalf("unexpected repository coordinator id: %s", captured.CoordinatorID)
	}
	if captured.ServiceTypeID != 2 {
		t.Fatalf("unexpected repository service_type_id: %d", captured.ServiceTypeID)
	}
	if captured.AssignedBy == nil || captured.AssignedBy.String() != actorID {
		t.Fatalf("expected assigned_by %s, got %#v", actorID, captured.AssignedBy)
	}
}

func TestListAssigned_CoordinatorUsesOwnUserID(t *testing.T) {
	actorID := uuid.New()
	var captured repository.ListAssignedInput

	svc := NewService(assignmentRepoMock{
		listAssignedFn: func(ctx context.Context, in repository.ListAssignedInput) (repository.ListResult, error) {
			captured = in
			return repository.ListResult{
				Items: []repository.ListItem{},
				Total: 0,
				Page:  in.Page,
				Limit: in.Limit,
			}, nil
		},
	}, nil)

	_, err := svc.ListAssigned(context.Background(), ListAssignedPreturnosInput{
		Page:        1,
		Limit:       20,
		ActorUserID: actorID.String(),
		ActorRoles:  []string{"COORDINADOR"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.AssignedCoordinatorID == nil || *captured.AssignedCoordinatorID != actorID {
		t.Fatalf("expected assigned coordinator id %s, got %#v", actorID, captured.AssignedCoordinatorID)
	}
}

func TestListAssigned_AdminRolesCanViewAllAssigned(t *testing.T) {
	var captured repository.ListAssignedInput

	svc := NewService(assignmentRepoMock{
		listAssignedFn: func(ctx context.Context, in repository.ListAssignedInput) (repository.ListResult, error) {
			captured = in
			return repository.ListResult{
				Items: []repository.ListItem{},
				Total: 0,
				Page:  in.Page,
				Limit: in.Limit,
			}, nil
		},
	}, nil)

	_, err := svc.ListAssigned(context.Background(), ListAssignedPreturnosInput{
		Page:        1,
		Limit:       20,
		ActorUserID: uuid.NewString(),
		ActorRoles:  []string{"SECRETARIA"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.AssignedCoordinatorID != nil {
		t.Fatalf("expected nil assigned coordinator filter, got %#v", captured.AssignedCoordinatorID)
	}
}

func TestListAssigned_ForbidsUnsupportedRole(t *testing.T) {
	svc := NewService(assignmentRepoMock{}, nil)

	_, err := svc.ListAssigned(context.Background(), ListAssignedPreturnosInput{
		Page:        1,
		Limit:       20,
		ActorUserID: uuid.NewString(),
		ActorRoles:  []string{"ESTUDIANTE"},
	})
	if !errors.Is(err, ErrAssignmentForbidden) {
		t.Fatalf("expected ErrAssignmentForbidden, got %v", err)
	}
}
