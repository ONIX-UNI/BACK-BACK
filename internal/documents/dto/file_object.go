package dto

import "time"

type FileObject struct {
	ID           string    `json:"id"`
	StorageKey   string    `json:"storage_key"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type,omitempty"`
	SizeBytes    int64     `json:"size_bytes"`
	Sha256       string    `json:"sha256,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateFileObjectRequest struct {
	StorageKey   string  `json:"storage_key" validate:"required"`
	OriginalName string  `json:"original_name" validate:"required"`
	MimeType     *string `json:"mime_type"`
	SizeBytes    int64   `json:"size_bytes" validate:"required"`
	Sha256       *string `json:"sha256"`
}

type UpdateFileObjectRequest struct {
	OriginalName *string `json:"original_name,omitempty"`
}

type GetFileObjectByIDRequest struct {
	ID string `params:"id" validate:"required,uuid"`
}

type GetFileObjectByNameRequest struct {
	OriginalName string `query:"name" validate:"required"`
}

type DeleteFileObjectRequest struct {
	ID string `params:"id" validate:"required,uuid"`
}
