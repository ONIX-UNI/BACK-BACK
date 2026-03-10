package service

import (
	"context"
	"mime/multipart"
	"time"

	documentsdto "github.com/DuvanRozoParra/sicou/internal/documents/dto"
)

type CreateAsesoriaInput struct {
	AcceptsDataProcessing  string
	ConsultationDate       string
	FullName               string
	DocumentType           string
	OtherDocumentType      string
	DocumentNumber         string
	BirthDate              string
	Age                    string
	MaritalStatus          string
	OtherMaritalStatus     string
	Gender                 string
	Address                string
	HousingType            string
	OtherHousingType       string
	SocioEconomicStratum   string
	SisbenCategory         string
	MobilePhone            string
	Email                  string
	PopulationType         string
	OtherPopulationType    string
	HeadOfHousehold        string
	Occupation             string
	EducationLevel         string
	OtherEducationLevel    string
	CaseDescription        string
	AuthorizesNotification string
	SubmittedAt            string
	IdentityDocument       *multipart.FileHeader
	UtilityBill            *multipart.FileHeader
	ActorUserID            string
	TimelineSource         string
}

const (
	TimelineSourcePublicForm = "formulario_publico"
	TimelineSourceInternal   = "lista_interna"
)

type CreateAsesoriaResult struct {
	ID             string
	PreturnoNumber string
	Status         string
	CreatedAt      time.Time
}

type ListPreturnosInput struct {
	Page  int
	Limit int
}

type ListAssignedPreturnosInput struct {
	Page        int
	Limit       int
	ActorUserID string
	ActorRoles  []string
}

type PreturnoSupport struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	URL         string `json:"url"`
	ViewURL     string `json:"viewUrl,omitempty"`
	DownloadURL string `json:"downloadUrl,omitempty"`
	Tipo        string `json:"tipo,omitempty"`
	Campo       string `json:"campo,omitempty"`
	Tamano      int64  `json:"tamano,omitempty"`
}

type PreturnoTimelineItem struct {
	ID      string `json:"id"`
	Fecha   string `json:"fecha"`
	Titulo  string `json:"titulo"`
	Detalle string `json:"detalle"`
	Usuario string `json:"usuario"`
	Source  string `json:"source,omitempty"`
}

type PreturnoListItem struct {
	ID                   string                 `json:"id"`
	Radicado             string                 `json:"radicado"`
	Turno                string                 `json:"turno"`
	Estado               string                 `json:"estado"`
	AutorizacionDatos    bool                   `json:"autorizacionDatos"`
	FechaConsulta        string                 `json:"fechaConsulta"`
	NombreCompleto       string                 `json:"nombreCompleto"`
	TipoDocumento        string                 `json:"tipoDocumento"`
	NumeroDocumento      string                 `json:"numeroDocumento"`
	FechaNacimiento      string                 `json:"fechaNacimiento"`
	Edad                 int                    `json:"edad"`
	EstadoCivil          string                 `json:"estadoCivil"`
	Genero               string                 `json:"genero"`
	Direccion            string                 `json:"direccion"`
	TipoVivienda         string                 `json:"tipoVivienda"`
	Estrato              *int                   `json:"estrato"`
	SisbenCategoria      string                 `json:"sisbenCategoria"`
	Telefono             string                 `json:"telefono"`
	CorreoElectronico    string                 `json:"correoElectronico"`
	TipoPoblacion        string                 `json:"tipoPoblacion"`
	CabezaHogar          *bool                  `json:"cabezaHogar"`
	Ocupacion            string                 `json:"ocupacion"`
	NivelEstudio         string                 `json:"nivelEstudio"`
	Relato               string                 `json:"relato"`
	Soportes             []PreturnoSupport      `json:"soportes"`
	AutorizaNotificacion bool                   `json:"autorizaNotificacion"`
	Timeline             []PreturnoTimelineItem `json:"timeline"`
	Citizen              string                 `json:"citizen"`
	Cedula               string                 `json:"cedula"`
}

type ListPreturnosResult struct {
	Items []PreturnoListItem `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}

type ValidationError struct {
	Fields map[string]string `json:"fields"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

type FileObjectService interface {
	Create(ctx context.Context, file multipart.File, size int64, filename string, contentType string) (*documentsdto.FileObject, error)
	Delete(ctx context.Context, id string) error
}
