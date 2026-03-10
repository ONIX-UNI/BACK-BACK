package dto

import (
	"time"

	"github.com/google/uuid"
)

type CaseFile struct {
	ID uuid.UUID `json:"id" db:"id"`

	CitizenID  uuid.UUID  `json:"citizen_id" db:"citizen_id"`
	PreturnoID uuid.UUID  `json:"preturno_id" db:"preturno_id"`
	TurnID     *uuid.UUID `json:"turn_id,omitempty" db:"turn_id"`

	ServiceTypeID int16  `json:"service_type_id" db:"service_type_id"`
	LegalAreaID   *int16 `json:"legal_area_id,omitempty" db:"legal_area_id"`

	Status string `json:"status" db:"status"`

	CurrentResponsible *uuid.UUID `json:"current_responsible,omitempty" db:"current_responsible"`
	SupervisorUser     *uuid.UUID `json:"supervisor_user,omitempty" db:"supervisor_user"`

	OpenedAt   time.Time  `json:"opened_at" db:"opened_at"`
	ClosedAt   *time.Time `json:"closed_at,omitempty" db:"closed_at"`
	CloseNotes *string    `json:"close_notes,omitempty" db:"close_notes"`

	CreatedBy *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty" db:"updated_by"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCaseFileRequest struct {
	CitizenID  uuid.UUID `json:"citizen_id" validate:"required"`
	PreturnoID uuid.UUID `json:"preturno_id" validate:"required"`

	ServiceTypeID int16  `json:"service_type_id" validate:"required"`
	LegalAreaID   *int16 `json:"legal_area_id,omitempty"`

	Status string `json:"status,omitempty"`

	CurrentResponsible *uuid.UUID `json:"current_responsible,omitempty"`
	SupervisorUser     *uuid.UUID `json:"supervisor_user,omitempty"`
}

type UpdateCaseFileRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`

	TurnID *uuid.UUID `json:"turn_id,omitempty"`

	ServiceTypeID *int16 `json:"service_type_id,omitempty"`
	LegalAreaID   *int16 `json:"legal_area_id,omitempty"`

	Status *string `json:"status,omitempty"`

	CurrentResponsible *uuid.UUID `json:"current_responsible,omitempty"`
	SupervisorUser     *uuid.UUID `json:"supervisor_user,omitempty"`

	ClosedAt   *time.Time `json:"closed_at,omitempty"`
	CloseNotes *string    `json:"close_notes,omitempty"`
}

type GetCaseFileByIDRequest struct {
	ID uuid.UUID `params:"id" validate:"required,uuid"`
}

type DeleteCaseFileRequest struct {
	ID uuid.UUID `params:"id" validate:"required,uuid"`
}

type ExpedienteResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Initials         string    `json:"initials"`
	Status           string    `json:"status"` // Activo | En revisión | Cerrado
	Cedula           string    `json:"cedula"`
	Area             string    `json:"area"`
	Tipo             string    `json:"tipo"`
	Creado           string    `json:"creado"`
	AttachmentsCount int       `json:"attachmentsCount"`
	CommentsCount    int       `json:"commentsCount"`
	Progress         int       `json:"progress"`
	CloseState       string    `json:"closeState"` // Bloqueado | Listo para cerrar
}

type CaseListRow struct {
	ID             uuid.UUID
	Status         string
	OpenedAt       time.Time
	FullName       string
	DocumentNumber string
	LegalArea      *string
	ServiceType    *string
}

type CaseListItem struct {
	ExpedienteCode string `json:"expediente_code"`

	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	OpenedAt time.Time  `json:"opened_at"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`

	SupervisorID   *uuid.UUID `json:"supervisor_id"`
	SupervisorName *string    `json:"supervisor_name,omitempty"`

	FullName       string `json:"full_name"`
	DocumentNumber string `json:"document_number"`

	LegalArea   *string `json:"legal_area,omitempty"`
	ServiceType *string `json:"service_type,omitempty"`

	ResponsibleName *string `json:"responsible_name,omitempty"`

	Consecutive *int       `json:"consecutive,omitempty"`
	TurnDate    *time.Time `json:"turn_date,omitempty"`

	Folders *string `json:"folders,omitempty"`

	AttachmentsCount int `json:"attachments_count"`
	CommentsCount    int `json:"comments_count"`
	Progress         int `json:"progress"`
}

type CaseOtpRequest struct {
	Email string `json:"email"`
}

type CaseOtpResponse struct {
	Message string `json:"message"`
}

type VerifyCaseFileOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,numeric,len=6"`
}

type CreateDocumentRequest struct {
	FileID         uuid.UUID  `json:"file_id" validate:"required"`
	DocumentKindID int16      `json:"document_kind_id" validate:"required"`
	PreturnoID     *uuid.UUID `json:"preturno_id,omitempty"`
	CaseID         *uuid.UUID `json:"case_id,omitempty"`
	CaseEventID    *uuid.UUID `json:"case_event_id,omitempty"`
	FolderID       *uuid.UUID `json:"folder_id,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
}

type DocumentResponse struct {
	ID             uuid.UUID  `json:"id"`
	FileID         uuid.UUID  `json:"file_id"`
	DocumentKindID int16      `json:"document_kind_id"`
	PreturnoID     *uuid.UUID `json:"preturno_id,omitempty"`
	CaseID         *uuid.UUID `json:"case_id,omitempty"`
	CaseEventID    *uuid.UUID `json:"case_event_id,omitempty"`
	Version        int        `json:"version"`
	IsCurrent      bool       `json:"is_current"`
	Notes          *string    `json:"notes,omitempty"`
	UploadedBy     *uuid.UUID `json:"uploaded_by,omitempty"`
	UploadedAt     time.Time  `json:"uploaded_at"`
}

type FolderDocumentResponse struct {
	FolderID   uuid.UUID  `json:"folder_id"`
	DocumentID uuid.UUID  `json:"document_id"`
	AddedAt    time.Time  `json:"added_at"`
	AddedBy    *uuid.UUID `json:"added_by,omitempty"`
}

type FolderResponse struct {
	ID uuid.UUID `json:"id"`

	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
