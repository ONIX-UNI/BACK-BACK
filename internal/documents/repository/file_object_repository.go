package repository

import (
	"context"
	"errors"

	"github.com/DuvanRozoParra/sicou/internal/documents/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository struct to handle database operations.
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new instance of the repository.
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create inserts a new file object into the database and returns it.
func (r *PostgresRepository) Create(ctx context.Context, file dto.FileObject) (*dto.FileObject, error) {

	query := `INSERT INTO sicou.file_object
		(id, storage_key, original_name, mime_type, size_bytes, sha256, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, storage_key, original_name, mime_type, size_bytes, sha256, created_at`
	row := r.db.QueryRow(ctx, query,
		file.ID,
		file.StorageKey,
		file.OriginalName,
		file.MimeType,
		file.SizeBytes,
		file.Sha256,
		file.CreatedAt)
	var result dto.FileObject

	err := row.Scan(
		&result.ID,
		&result.StorageKey,
		&result.OriginalName,
		&result.MimeType,
		&result.SizeBytes,
		&result.Sha256,
		&result.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *PostgresRepository) UpdateFull(
	ctx context.Context,
	file dto.FileObject,
) (*dto.FileObject, error) {

	if file.ID == "" {
		return nil, errors.New("file ID is required")
	}

	query := `
		UPDATE sicou.file_object
		SET storage_key=$1,
			original_name=$2,
			mime_type=$3,
			size_bytes=$4,
			sha256=$5
		WHERE id=$6
		RETURNING id, storage_key, original_name, mime_type, size_bytes, sha256, created_at
	`

	row := r.db.QueryRow(ctx, query,
		file.StorageKey,
		file.OriginalName,
		file.MimeType,
		file.SizeBytes,
		file.Sha256,
		file.ID,
	)

	var result dto.FileObject

	err := row.Scan(
		&result.ID,
		&result.StorageKey,
		&result.OriginalName,
		&result.MimeType,
		&result.SizeBytes,
		&result.Sha256,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes an existing file object from the database by its ID.
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("file ID is required")
	}

	query := `DELETE FROM sicou.file_object WHERE id=$1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("file not found")
	}

	return nil
}

// GetAll retrieves all file objects from the database.
func (r *PostgresRepository) GetAll(ctx context.Context) ([]dto.FileObject, error) {

	query := `
		SELECT id, storage_key, original_name, mime_type, size_bytes, sha256, created_at
		FROM sicou.file_object
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []dto.FileObject

	for rows.Next() {
		var file dto.FileObject
		err := rows.Scan(
			&file.ID,
			&file.StorageKey,
			&file.OriginalName,
			&file.MimeType,
			&file.SizeBytes,
			&file.Sha256,
			&file.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

// GetById retrieves a file object from the database by its ID.
func (r *PostgresRepository) GetById(ctx context.Context, id string) (*dto.FileObject, error) {

	if id == "" {
		return nil, errors.New("file ID is required")
	}

	query := `
		SELECT id, storage_key, original_name, mime_type, size_bytes, sha256, created_at
		FROM sicou.file_object
		WHERE id=$1
	`

	row := r.db.QueryRow(ctx, query, id)

	var file dto.FileObject

	err := row.Scan(
		&file.ID,
		&file.StorageKey,
		&file.OriginalName,
		&file.MimeType,
		&file.SizeBytes,
		&file.Sha256,
		&file.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &file, nil
}

func (r *PostgresRepository) GetByName(
	ctx context.Context,
	name string,
) (*dto.FileObject, error) {

	if name == "" {
		return nil, errors.New("original name is required")
	}

	query := `
		SELECT id, storage_key, original_name, mime_type, size_bytes, sha256, created_at
		FROM sicou.file_object
		WHERE original_name=$1
		LIMIT 1
	`

	row := r.db.QueryRow(ctx, query, name)

	var file dto.FileObject

	err := row.Scan(
		&file.ID,
		&file.StorageKey,
		&file.OriginalName,
		&file.MimeType,
		&file.SizeBytes,
		&file.Sha256,
		&file.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &file, nil
}

func (r *PostgresRepository) List(
	ctx context.Context,
	limit, offset int,
) ([]dto.FileObject, error) {

	query := `
		SELECT id, storage_key, original_name, mime_type, size_bytes, sha256, created_at
		FROM sicou.file_object
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []dto.FileObject

	for rows.Next() {
		var file dto.FileObject
		err := rows.Scan(
			&file.ID,
			&file.StorageKey,
			&file.OriginalName,
			&file.MimeType,
			&file.SizeBytes,
			&file.Sha256,
			&file.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}
