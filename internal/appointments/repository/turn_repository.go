package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type STurn struct {
	db *pgxpool.Pool
}

func NewTurnRepository(db *pgxpool.Pool) *STurn {
	return &STurn{db: db}
}

func (r *STurn) Create(ctx context.Context, req dto.CreateTurnRequest) (*dto.Turn, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var (
		citizenID     uuid.UUID
		serviceTypeID int16
		legalAreaID   *int16
		coordinatorID *uuid.UUID
		currentStatus string
	)
	err = tx.QueryRow(ctx, `
		SELECT
			p.citizen_id,
			COALESCE(p.service_type_id, 0) AS service_type_id,
			st.legal_area_id,
			p.assigned_coordinator_id,
			p.status
		FROM sicou.preturno p
		LEFT JOIN sicou.service_type st ON st.id = p.service_type_id
		WHERE p.id = $1
	`, req.PreturnoID).Scan(
		&citizenID,
		&serviceTypeID,
		&legalAreaID,
		&coordinatorID,
		&currentStatus,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("preturno not found")
		}
		return nil, err
	}
	if serviceTypeID <= 0 {
		return nil, errors.New("preturno has no service_type assigned")
	}

	query := `
		INSERT INTO sicou.turn (
			preturno_id,
			turn_date,
			consecutive,
			priority,
			created_by
		)
		VALUES (
			$1,
			COALESCE($2, now()::date),
			(
				SELECT COALESCE(MAX(consecutive), 0) + 1
				FROM sicou.turn
				WHERE turn_date = COALESCE($2, now()::date)
			),
			COALESCE($3, 0),
			$4
		)
		RETURNING
			id, preturno_id, turn_date, consecutive,
			status, priority,
			escritorio_id,
			created_by, created_at,
			called_at, attended_at, finished_at
	`

	var turn dto.Turn
	err = tx.QueryRow(ctx, query,
		req.PreturnoID,
		req.TurnDate,
		req.Priority,
		req.CreatedBy,
	).Scan(
		&turn.ID,
		&turn.PreturnoID,
		&turn.TurnDate,
		&turn.Consecutive,
		&turn.Status,
		&turn.Priority,
		&turn.EscritorioID,
		&turn.CreatedBy,
		&turn.CreatedAt,
		&turn.CalledAt,
		&turn.AttendedAt,
		&turn.FinishedAt,
	)

	if err != nil {
		return nil, err
	}

	assignedStudentID := req.StudentID
	if assignedStudentID == nil {
		assignedStudentID = req.AssignedTo
	}

	currentResponsible := assignedStudentID
	if currentResponsible == nil {
		currentResponsible = req.CreatedBy
	}

	supervisorUser := coordinatorID
	if supervisorUser == nil {
		supervisorUser = req.CreatedBy
	}

	if assignedStudentID != nil && req.CreatedBy != nil {
		_, err = tx.Exec(ctx, `
			INSERT INTO sicou.turn_assignment (
				turn_id,
				stage,
				method,
				assigned_to,
				assigned_by,
				reason
			)
			VALUES ($1, 'DEFINITIVO', 'MANUAL', $2, $3, $4)
		`,
			turn.ID,
			assignedStudentID,
			req.CreatedBy,
			"Asignacion de estudiante al crear turno",
		)
		if err != nil {
			return nil, err
		}
	}

	var (
		caseID      uuid.UUID
		caseCreated bool
	)
	err = tx.QueryRow(ctx, `
		INSERT INTO sicou.case_file (
			citizen_id,
			preturno_id,
			turn_id,
			service_type_id,
			legal_area_id,
			status,
			current_responsible,
			supervisor_user,
			created_by,
			updated_by
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			'ABIERTO',
			$6,
			$7,
			$8,
			$8
		)
		ON CONFLICT (preturno_id)
		DO UPDATE SET
			turn_id = EXCLUDED.turn_id,
			service_type_id = EXCLUDED.service_type_id,
			legal_area_id = EXCLUDED.legal_area_id,
			current_responsible = COALESCE(EXCLUDED.current_responsible, sicou.case_file.current_responsible),
			supervisor_user = COALESCE(EXCLUDED.supervisor_user, sicou.case_file.supervisor_user),
			updated_by = COALESCE(EXCLUDED.updated_by, sicou.case_file.updated_by),
			updated_at = now()
		RETURNING id, (xmax = 0) AS inserted
	`,
		citizenID,
		req.PreturnoID,
		turn.ID,
		serviceTypeID,
		legalAreaID,
		currentResponsible,
		supervisorUser,
		req.CreatedBy,
	).Scan(&caseID, &caseCreated)
	if err != nil {
		return nil, err
	}

	caseEventType := "TURN_LINKED"
	caseEventTitle := "Turno vinculado al caso"
	caseEventNote := "Se vinculo el turno al expediente desde preturno"
	if caseCreated {
		caseEventType = "CREATED"
		caseEventTitle = "Caso creado desde preturno"
		caseEventNote = "Creacion automatica del caso al generar el turno"
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO sicou.case_event (
			case_id,
			event_type,
			title,
			notes,
			created_by
		)
		VALUES ($1, $2, $3, $4, $5)
	`,
		caseID,
		caseEventType,
		caseEventTitle,
		caseEventNote,
		req.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(strings.TrimSpace(currentStatus), "EN_TURNO") {
		_, err = tx.Exec(ctx, `
			UPDATE sicou.preturno
			SET status = 'EN_TURNO',
				updated_by = COALESCE($2, updated_by),
				updated_at = now()
			WHERE id = $1
		`, req.PreturnoID, req.CreatedBy)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &turn, nil
}
func (r *STurn) List(ctx context.Context) ([]dto.Turn, error) {
	query := `
		SELECT 
			id, preturno_id, turn_date, consecutive,
			status, priority,
			escritorio_id,
			created_by, created_at,
			called_at, attended_at, finished_at
		FROM sicou.turn
		ORDER BY turn_date DESC, consecutive DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turns []dto.Turn

	for rows.Next() {
		var turn dto.Turn
		err := rows.Scan(
			&turn.ID,
			&turn.PreturnoID,
			&turn.TurnDate,
			&turn.Consecutive,
			&turn.Status,
			&turn.Priority,
			&turn.EscritorioID,
			&turn.CreatedBy,
			&turn.CreatedAt,
			&turn.CalledAt,
			&turn.AttendedAt,
			&turn.FinishedAt,
		)
		if err != nil {
			return nil, err
		}

		turns = append(turns, turn)
	}

	return turns, rows.Err()
}
func (r *STurn) GetById(ctx context.Context, id uuid.UUID) (*dto.Turn, error) {
	query := `
		SELECT 
			id, preturno_id, turn_date, consecutive,
			status, priority,
			created_by, created_at,
			called_at, attended_at, finished_at
		FROM sicou.turn
		WHERE id = $1
	`

	var turn dto.Turn

	err := r.db.QueryRow(ctx, query, id).Scan(
		&turn.ID,
		&turn.PreturnoID,
		&turn.TurnDate,
		&turn.Consecutive,
		&turn.Status,
		&turn.Priority,
		&turn.CreatedBy,
		&turn.CreatedAt,
		&turn.CalledAt,
		&turn.AttendedAt,
		&turn.FinishedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &turn, nil
}
func (r *STurn) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTurnRequest) (*dto.Turn, error) {
	setParts := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *req.Status)
		argPos++
	}

	if req.Priority != nil {
		setParts = append(setParts, fmt.Sprintf("priority = $%d", argPos))
		args = append(args, *req.Priority)
		argPos++
	}

	if req.CalledAt != nil {
		setParts = append(setParts, fmt.Sprintf("called_at = $%d", argPos))
		args = append(args, *req.CalledAt)
		argPos++
	}

	if req.AttendedAt != nil {
		setParts = append(setParts, fmt.Sprintf("attended_at = $%d", argPos))
		args = append(args, *req.AttendedAt)
		argPos++
	}

	if req.FinishedAt != nil {
		setParts = append(setParts, fmt.Sprintf("finished_at = $%d", argPos))
		args = append(args, *req.FinishedAt)
		argPos++
	}

	if len(setParts) == 0 {
		return r.GetById(ctx, id)
	}

	query := fmt.Sprintf(`
		UPDATE sicou.turn
		SET %s
		WHERE id = $%d
		RETURNING 
			id, preturno_id, turn_date, consecutive,
			status, priority,
			created_by, created_at,
			called_at, attended_at, finished_at
	`, strings.Join(setParts, ", "), argPos)

	args = append(args, id)

	var turn dto.Turn
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&turn.ID,
		&turn.PreturnoID,
		&turn.TurnDate,
		&turn.Consecutive,
		&turn.Status,
		&turn.Priority,
		&turn.CreatedBy,
		&turn.CreatedAt,
		&turn.CalledAt,
		&turn.AttendedAt,
		&turn.FinishedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &turn, nil
}
func (r *STurn) Delete(ctx context.Context, id uuid.UUID) (*dto.Turn, error) {
	query := `
		DELETE FROM sicou.turn
		WHERE id = $1
		RETURNING 
			id, preturno_id, turn_date, consecutive,
			status, priority,
			created_by, created_at,
			called_at, attended_at, finished_at
	`

	var turn dto.Turn

	err := r.db.QueryRow(ctx, query, id).Scan(
		&turn.ID,
		&turn.PreturnoID,
		&turn.TurnDate,
		&turn.Consecutive,
		&turn.Status,
		&turn.Priority,
		&turn.CreatedBy,
		&turn.CreatedAt,
		&turn.CalledAt,
		&turn.AttendedAt,
		&turn.FinishedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &turn, nil
}

func (r *STurn) SetTurnDesktop(ctx context.Context, req dto.SetTurnDesktopRequest) (*dto.SetTurnEscritorioResponse, error) {

	if req.TurnID == uuid.Nil {
		return nil, fmt.Errorf("invalid turn_id")
	}

	if req.EscritorioID == uuid.Nil {
		return nil, fmt.Errorf("invalid escritorio_id")
	}

	query := `
		SELECT sicou.fn_turn_set_escritorio($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query,
		req.TurnID,
		req.EscritorioID,
		req.Reason,
	)

	if err != nil {
		return nil, err
	}

	return &dto.SetTurnEscritorioResponse{
		Success: true,
		Message: "desktop assigned successfully",
	}, nil
}
