package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	StatusPendientePreturno = "PENDIENTE_VALIDACION"
	StatusAsignadoPreturno  = "ASIGNADO"
)

type AttachmentRecord struct {
	FileObjectID string
	FieldName    string
	OriginalName string
	MimeType     string
	SizeBytes    int64
	URL          string
	StorageKey   string
}

type CreateRecord struct {
	ID                     string
	Payload                string
	ConsultationDate       time.Time
	SubmittedAt            time.Time
	AcceptsDataProcessing  bool
	AuthorizesNotification bool
	NotificationEmail      string
	CitizenName            string
	HeadOfHousehold        *bool
	CreatedBy              *uuid.UUID
	EventSource            string
	Attachments            []AttachmentRecord
}

type CreateResult struct {
	ID             string
	PreturnoNumber string
	Status         string
	CreatedAt      time.Time
}

type ListInput struct {
	Page  int
	Limit int
}

type ListAssignedInput struct {
	Page                  int
	Limit                 int
	AssignedCoordinatorID *uuid.UUID
}

type IntakePayload struct {
	AcceptsDataProcessing  string `json:"acceptsDataProcessing"`
	ConsultationDate       string `json:"consultationDate"`
	FullName               string `json:"fullName"`
	DocumentType           string `json:"documentType"`
	OtherDocumentType      string `json:"otherDocumentType"`
	DocumentNumber         string `json:"documentNumber"`
	BirthDate              string `json:"birthDate"`
	Age                    string `json:"age"`
	MaritalStatus          string `json:"maritalStatus"`
	OtherMaritalStatus     string `json:"otherMaritalStatus"`
	Gender                 string `json:"gender"`
	Address                string `json:"address"`
	HousingType            string `json:"housingType"`
	OtherHousingType       string `json:"otherHousingType"`
	SocioEconomicStratum   string `json:"socioEconomicStratum"`
	SisbenCategory         string `json:"sisbenCategory"`
	MobilePhone            string `json:"mobilePhone"`
	Email                  string `json:"email"`
	PopulationType         string `json:"populationType"`
	OtherPopulationType    string `json:"otherPopulationType"`
	HeadOfHousehold        string `json:"headOfHousehold"`
	Occupation             string `json:"occupation"`
	EducationLevel         string `json:"educationLevel"`
	OtherEducationLevel    string `json:"otherEducationLevel"`
	CaseDescription        string `json:"caseDescription"`
	AuthorizesNotification string `json:"authorizesNotification"`
	SubmittedAt            string `json:"submittedAt"`
}

type ListAttachment struct {
	ID           string
	FileObjectID string
	FieldName    string
	OriginalName string
	MimeType     string
	SizeBytes    int64
	URL          string
	CreatedAt    time.Time
}

type ListTimelineItem struct {
	ID        string
	Title     string
	Detail    string
	User      string
	Source    string
	CreatedAt time.Time
}

type ListItem struct {
	ID               string
	PreturnoNumber   string
	Status           string
	CreatedAt        time.Time
	ConsultationDate time.Time
	Payload          IntakePayload
	Attachments      []ListAttachment
	Timeline         []ListTimelineItem
}

type ListResult struct {
	Items []ListItem
	Total int64
	Page  int
	Limit int
}

type AssignPreturnoInput struct {
	PreturnoID    string
	CoordinatorID string
	ServiceTypeID int16
	Observations  string
	AssignedBy    *uuid.UUID
}

type AssignmentTimelineEvent struct {
	ID        string
	Title     string
	Detail    string
	CreatedAt time.Time
}

type AssignPreturnoResult struct {
	ID                    string
	Status                string
	AssignedCoordinatorID string
	ServiceTypeID         int16
	TimelineEvent         AssignmentTimelineEvent
}

type AssignmentCoordinatorOption struct {
	ID          string
	DisplayName string
	Email       string
}

type AssignmentServiceTypeOption struct {
	ID   int16
	Code string
	Name string
}

type AssignmentOptionsResult struct {
	Coordinators []AssignmentCoordinatorOption
	ServiceTypes []AssignmentServiceTypeOption
}

type Repository interface {
	Create(ctx context.Context, in CreateRecord) (CreateResult, error)
	List(ctx context.Context, in ListInput) (ListResult, error)
	ListAssigned(ctx context.Context, in ListAssignedInput) (ListResult, error)
	AssignPreturno(ctx context.Context, in AssignPreturnoInput) (AssignPreturnoResult, error)
	AssignmentOptions(ctx context.Context) (AssignmentOptionsResult, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}
