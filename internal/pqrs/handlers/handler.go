package handlers

import (
	"errors"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/pqrs/dto"
	"github.com/DuvanRozoParra/sicou/internal/pqrs/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
	Detail  string            `json:"detail,omitempty"`
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *fiber.Ctx) error {
	contentType := strings.ToLower(strings.TrimSpace(c.Get(fiber.HeaderContentType)))
	if !strings.Contains(contentType, "multipart/form-data") {
		log.Printf("pqrs create rejected: invalid content-type=%q", contentType)
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(ErrorResponse{
			Error:   "UNSUPPORTED_MEDIA_TYPE",
			Message: "content-type must be multipart/form-data",
			Detail:  contentType,
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("pqrs create rejected: invalid multipart body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "INVALID_MULTIPART_BODY",
			Message: "invalid multipart/form-data body",
			Detail:  err.Error(),
		})
	}

	result, err := h.service.Create(c.Context(), service.CreatePQRSInput{
		QueryType:                 formValue(form, "queryType"),
		PersonType:                formValue(form, "personType"),
		DocumentType:              formValue(form, "documentType"),
		DocumentOrTaxID:           formValue(form, "documentOrTaxId"),
		FirstName:                 formValue(form, "firstName"),
		MiddleName:                formValue(form, "middleName"),
		FirstLastName:             formValue(form, "firstLastName"),
		SecondLastName:            formValue(form, "secondLastName"),
		Gender:                    formValue(form, "gender"),
		Address:                   formValue(form, "address"),
		NeighborhoodArea:          formValue(form, "neighborhoodArea"),
		AllowsElectronicResponse:  formValue(form, "allowsElectronicResponse"),
		Email:                     formValue(form, "email"),
		Phone:                     formValue(form, "phone"),
		PopulationGroup:           formValue(form, "populationGroup"),
		OtherPopulationGroup:      formValue(form, "otherPopulationGroup"),
		RequestDescription:        formValue(form, "requestDescription"),
		ResponseChannel:           formValue(form, "responseChannel"),
		RequestType:               formValue(form, "requestType"),
		RequestAgainstStudent:     formValue(form, "requestAgainstStudent"),
		ResponsibleStudentName:    formValue(form, "responsibleStudentName"),
		ResponsibleStudentProgram: formValue(form, "responsibleStudentProgram"),
		StudentCaseDescription:    formValue(form, "studentCaseDescription"),
		AcceptsDataProcessing:     formValue(form, "acceptsDataProcessing"),
		SubmittedAt:               formValue(form, "submittedAt"),
		Attachments:               attachmentFiles(form),
	})
	if err != nil {
		var validationErr *service.ValidationError
		if errors.As(err, &validationErr) {
			log.Printf("pqrs validation failed: %+v", validationErr.Fields)
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "VALIDATION_ERROR",
				Message: validationErr.Error(),
				Fields:  validationErr.Fields,
			})
		}

		log.Printf("pqrs create failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "PQRS_CREATE_FAILED",
			Message: "failed to register pqrs",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CreatePQRSResponse{
		ID:        result.ID,
		Radicado:  result.Radicado,
		Estado:    result.Estado,
		CreatedAt: result.CreatedAt,
	})
}

func (h *Handler) List(c *fiber.Ctx) error {
	page := parsePositiveInt(c.Query("page"), 1)
	limit := parsePositiveInt(c.Query("limit"), 10)

	result, err := h.service.List(c.Context(), service.ListPQRSInput{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
		Estado: c.Query("estado"),
		Tipo:   c.Query("tipo"),
	})
	if err != nil {
		log.Printf("pqrs list failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "PQRS_LIST_FAILED",
			Message: "failed to list pqrs",
			Detail:  err.Error(),
		})
	}

	items := make([]dto.PQRSListItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, dto.PQRSListItem{
			ID:           item.ID,
			Radicado:     item.Radicado,
			Tipo:         item.Tipo,
			Asunto:       item.Asunto,
			Ciudadano:    item.Ciudadano,
			Correo:       item.Correo,
			Estado:       item.Estado,
			Responsable:  item.Responsable,
			FechaIngreso: item.FechaIngreso,
			FechaLimite:  item.FechaLimite,
		})
	}

	return c.JSON(dto.ListPQRSResponse{
		Data: items,
		Meta: dto.PQRSListMeta{
			Page:  result.Page,
			Limit: result.Limit,
			Total: result.Total,
		},
	})
}

func formValue(form *multipart.Form, key string) string {
	if form == nil {
		return ""
	}

	values := form.Value[key]
	if len(values) == 0 {
		return ""
	}

	return strings.TrimSpace(values[0])
}

func attachmentFiles(form *multipart.Form) []*multipart.FileHeader {
	if form == nil {
		return nil
	}

	files := make([]*multipart.FileHeader, 0)

	for _, key := range []string{"attachments", "attachments[]", "adjuntos", "adjuntos[]"} {
		if selected, ok := form.File[key]; ok {
			files = append(files, selected...)
		}
	}

	for key, selected := range form.File {
		if strings.HasPrefix(key, "attachments[") || strings.HasPrefix(key, "adjuntos[") {
			files = append(files, selected...)
		}
	}

	return files
}

func parsePositiveInt(value string, defaultValue int) int {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return defaultValue
	}

	return parsed
}
