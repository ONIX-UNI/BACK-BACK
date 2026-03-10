package repository

import (
	"context"
	"errors"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/citizens/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SCitizens struct {
	db *pgxpool.Pool
}

func NewCitizensRepository(db *pgxpool.Pool) *SCitizens {
	return &SCitizens{db: db}
}

func (r *SCitizens) Create(ctx context.Context, req dto.CreateCitizenRequest) (*dto.Citizen, error) {
	query := `
		INSERT INTO sicou.citizen (
			document_type_id,
			document_number,
			full_name,
			birth_date,
			phone_mobile,
			email,
			address
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING 
			id, document_type_id, document_number, full_name,
			birth_date, phone_mobile, email, address,
			created_at, updated_at
	`

	var c dto.Citizen

	err := r.db.QueryRow(ctx, query,
		req.DocumentTypeID,
		req.DocumentNumber,
		req.FullName,
		req.BirthDate,
		req.PhoneMobile,
		req.Email,
		req.Address,
	).Scan(
		&c.ID,
		&c.DocumentTypeID,
		&c.DocumentNumber,
		&c.FullName,
		&c.BirthDate,
		&c.PhoneMobile,
		&c.Email,
		&c.Address,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
func (r *SCitizens) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCitizenRequest) (*dto.Citizen, error) {
	query := `
		UPDATE sicou.citizen
		SET
			document_type_id = $1,
			document_number = $2,
			full_name = $3,
			birth_date = $4,
			phone_mobile = $5,
			email = $6,
			address = $7,
			updated_at = now()
		WHERE id = $8 AND deleted_at IS NULL
		RETURNING 
			id, document_type_id, document_number, full_name,
			birth_date, phone_mobile, email, address,
			created_at, updated_at, deleted_at
	`

	var c dto.Citizen

	err := r.db.QueryRow(ctx, query,
		req.DocumentTypeID,
		req.DocumentNumber,
		req.FullName,
		req.BirthDate,
		req.PhoneMobile,
		req.Email,
		req.Address,
		id,
	).Scan(
		&c.ID,
		&c.DocumentTypeID,
		&c.DocumentNumber,
		&c.FullName,
		&c.BirthDate,
		&c.PhoneMobile,
		&c.Email,
		&c.Address,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}
func (r *SCitizens) List(ctx context.Context, limit, offset int) ([]dto.Citizen, error) {
	query := `
		SELECT 
			id, document_type_id, document_number, full_name,
			birth_date, phone_mobile, email, address,
			created_at, updated_at, deleted_at
		FROM sicou.citizen
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var citizens []dto.Citizen

	for rows.Next() {
		var c dto.Citizen

		err := rows.Scan(
			&c.ID,
			&c.DocumentTypeID,
			&c.DocumentNumber,
			&c.FullName,
			&c.BirthDate,
			&c.PhoneMobile,
			&c.Email,
			&c.Address,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		citizens = append(citizens, c)
	}

	return citizens, nil
}
func (r *SCitizens) GetById(ctx context.Context, id uuid.UUID) (*dto.Citizen, error) {
	query := `
		SELECT 
			id, document_type_id, document_number, full_name,
			birth_date, phone_mobile, email, address,
			created_at, updated_at, deleted_at
		FROM sicou.citizen
		WHERE id = $1 AND deleted_at IS NULL
	`

	var c dto.Citizen

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.DocumentTypeID,
		&c.DocumentNumber,
		&c.FullName,
		&c.BirthDate,
		&c.PhoneMobile,
		&c.Email,
		&c.Address,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}
func (r *SCitizens) Delete(ctx context.Context, id uuid.UUID) (*dto.Citizen, error) {
	query := `
		UPDATE sicou.citizen
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING 
			id, document_type_id, document_number, full_name,
			birth_date, phone_mobile, email, address,
			created_at, updated_at, deleted_at
	`

	now := time.Now()

	var c dto.Citizen

	err := r.db.QueryRow(ctx, query, now, id).Scan(
		&c.ID,
		&c.DocumentTypeID,
		&c.DocumentNumber,
		&c.FullName,
		&c.BirthDate,
		&c.PhoneMobile,
		&c.Email,
		&c.Address,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}
