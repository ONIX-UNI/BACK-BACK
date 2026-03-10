package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

type StatsCounts struct {
	TurnosAtendidosHoyCurrent int64
	TurnosAtendidosHoyPrev    int64
	TurnosPendientesCurrent   int64
	TurnosPendientesPrev      int64
	CasosCreadosHoyCurrent    int64
	CasosCreadosHoyPrev       int64
	AtencionesEnCursoCurrent  int64
	AtencionesEnCursoPrev     int64
}

type PendingQueueRow struct {
	Consecutive int
	Citizen     string
	Motivo      string
	Canal       string
	CreatedAt   time.Time
	Priority    int
}

type DeadlineAlertRow struct {
	CaseID      string
	CaseCode    string
	Responsible string
	DueDate     time.Time
	TermStatus  string
}

type RecentActivityRow struct {
	CaseCode  string
	UserName  string
	EventType string
	CreatedAt time.Time
}

func (r *Repository) OverviewStats(ctx context.Context, day time.Time, loc *time.Location) (StatsCounts, error) {
	result := StatsCounts{}
	if r.db == nil {
		return result, nil
	}

	from := day.Format("2006-01-02")
	to := day.AddDate(0, 0, 1).Format("2006-01-02")
	prevFrom := day.AddDate(0, 0, -1).Format("2006-01-02")
	prevTo := day.Format("2006-01-02")
	tz := loc.String()

	const query = `
SELECT
  COUNT(*) FILTER (
    WHERE t.status = 'FINALIZADO'
      AND (t.finished_at AT TIME ZONE $5) >= $1::timestamp
      AND (t.finished_at AT TIME ZONE $5) < $2::timestamp
  ) AS attended_curr,
  COUNT(*) FILTER (
    WHERE t.status = 'FINALIZADO'
      AND (t.finished_at AT TIME ZONE $5) >= $3::timestamp
      AND (t.finished_at AT TIME ZONE $5) < $4::timestamp
  ) AS attended_prev,
  COUNT(*) FILTER (
    WHERE t.status IN ('EN_COLA','LLAMADO')
      AND (t.created_at AT TIME ZONE $5) < $2::timestamp
  ) AS pending_curr,
  COUNT(*) FILTER (
    WHERE t.status IN ('EN_COLA','LLAMADO')
      AND (t.created_at AT TIME ZONE $5) < $4::timestamp
  ) AS pending_prev,
  (SELECT COUNT(*)
    FROM sicou.case_file cf
    WHERE (cf.created_at AT TIME ZONE $5) >= $1::timestamp
      AND (cf.created_at AT TIME ZONE $5) < $2::timestamp
  ) AS cases_curr,
  (SELECT COUNT(*)
    FROM sicou.case_file cf
    WHERE (cf.created_at AT TIME ZONE $5) >= $3::timestamp
      AND (cf.created_at AT TIME ZONE $5) < $4::timestamp
  ) AS cases_prev,
  COUNT(*) FILTER (
    WHERE t.status = 'EN_ATENCION'
      AND (t.created_at AT TIME ZONE $5) < $2::timestamp
  ) AS in_progress_curr,
  COUNT(*) FILTER (
    WHERE t.status = 'EN_ATENCION'
      AND (t.created_at AT TIME ZONE $5) < $4::timestamp
  ) AS in_progress_prev
FROM sicou.turn t;
`

	err := r.db.QueryRow(ctx, query, from, to, prevFrom, prevTo, tz).Scan(
		&result.TurnosAtendidosHoyCurrent,
		&result.TurnosAtendidosHoyPrev,
		&result.TurnosPendientesCurrent,
		&result.TurnosPendientesPrev,
		&result.CasosCreadosHoyCurrent,
		&result.CasosCreadosHoyPrev,
		&result.AtencionesEnCursoCurrent,
		&result.AtencionesEnCursoPrev,
	)
	if err != nil {
		if isSchemaDriftError(err) {
			return StatsCounts{}, nil
		}
		return StatsCounts{}, err
	}

	return result, nil
}

func (r *Repository) PendingQueue(ctx context.Context, limit int) ([]PendingQueueRow, error) {
	if r.db == nil || limit <= 0 {
		return []PendingQueueRow{}, nil
	}

	const query = `
SELECT
  COALESCE(t.consecutive, 0) AS consecutive,
  COALESCE(NULLIF(TRIM(p.full_name_snapshot), ''), c.full_name, 'Sin nombre') AS citizen,
  COALESCE(NULLIF(TRIM(p.situation_story), ''), 'Sin descripcion') AS motivo,
  COALESCE(NULLIF(TRIM(ch.name), ''), 'No definido') AS canal,
  t.created_at,
  COALESCE(t.priority, 0) AS priority
FROM sicou.turn t
INNER JOIN sicou.preturno p ON p.id = t.preturno_id
LEFT JOIN sicou.citizen c ON c.id = p.citizen_id
LEFT JOIN sicou.catalog_channel ch ON ch.id = p.channel_id
WHERE t.status IN ('EN_COLA', 'LLAMADO')
ORDER BY t.priority DESC, t.created_at ASC
LIMIT $1;
`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		if isSchemaDriftError(err) {
			return []PendingQueueRow{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	result := make([]PendingQueueRow, 0, limit)
	for rows.Next() {
		var row PendingQueueRow
		if err := rows.Scan(
			&row.Consecutive,
			&row.Citizen,
			&row.Motivo,
			&row.Canal,
			&row.CreatedAt,
			&row.Priority,
		); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, rows.Err()
}

func (r *Repository) DeadlineAlerts(ctx context.Context, day time.Time, limit int) ([]DeadlineAlertRow, error) {
	if r.db == nil || limit <= 0 {
		return []DeadlineAlertRow{}, nil
	}

	from := day.Format("2006-01-02")

	const query = `
SELECT
  cf.id::text AS case_id,
  ('EXP-' || EXTRACT(YEAR FROM cf.opened_at)::int || '-' || LPAD(COALESCE(t.consecutive,0)::text, 4, '0')) AS case_code,
  COALESCE(NULLIF(TRIM(u.display_name), ''), 'Sin responsable') AS responsible,
  ct.due_date,
  ct.status
FROM sicou.case_term ct
INNER JOIN sicou.case_file cf ON cf.id = ct.case_id
LEFT JOIN sicou.turn t ON t.id = cf.turn_id
LEFT JOIN sicou.app_user u ON u.id = cf.current_responsible
WHERE ct.status IN ('VENCIDO','PROXIMO','ATIEMPO')
  AND ct.due_date >= $1::date - INTERVAL '15 days'
ORDER BY
  CASE ct.status
    WHEN 'VENCIDO' THEN 0
    WHEN 'PROXIMO' THEN 1
    ELSE 2
  END,
  ct.due_date ASC
LIMIT $2;
`

	rows, err := r.db.Query(ctx, query, from, limit)
	if err != nil {
		if isSchemaDriftError(err) {
			return []DeadlineAlertRow{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	result := make([]DeadlineAlertRow, 0, limit)
	for rows.Next() {
		var row DeadlineAlertRow
		if err := rows.Scan(&row.CaseID, &row.CaseCode, &row.Responsible, &row.DueDate, &row.TermStatus); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, rows.Err()
}

func (r *Repository) RecentActivity(ctx context.Context, limit int) ([]RecentActivityRow, error) {
	if r.db == nil || limit <= 0 {
		return []RecentActivityRow{}, nil
	}

	const query = `
SELECT
  ('EXP-' || EXTRACT(YEAR FROM cf.opened_at)::int || '-' || LPAD(COALESCE(t.consecutive,0)::text, 4, '0')) AS case_code,
  COALESCE(NULLIF(TRIM(u.display_name), ''), 'Sistema') AS user_name,
  COALESCE(NULLIF(TRIM(ce.event_type), ''), 'UPDATED') AS event_type,
  ce.created_at
FROM sicou.case_event ce
INNER JOIN sicou.case_file cf ON cf.id = ce.case_id
LEFT JOIN sicou.turn t ON t.id = cf.turn_id
LEFT JOIN sicou.app_user u ON u.id = ce.created_by
ORDER BY ce.created_at DESC
LIMIT $1;
`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		if isSchemaDriftError(err) {
			return []RecentActivityRow{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	result := make([]RecentActivityRow, 0, limit)
	for rows.Next() {
		var row RecentActivityRow
		if err := rows.Scan(&row.CaseCode, &row.UserName, &row.EventType, &row.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, rows.Err()
}

func isSchemaDriftError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}

	switch pgErr.Code {
	case "42P01", // undefined_table
		"42703": // undefined_column
		return true
	default:
		return false
	}
}
