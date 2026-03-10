package handlers

import (
	"errors"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/DuvanRozoParra/sicou/internal/appointments/service"
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
		log.Printf("form/legal-advise create rejected: invalid content-type=%q", contentType)
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(ErrorResponse{
			Error:   "UNSUPPORTED_MEDIA_TYPE",
			Message: "content-type must be multipart/form-data",
			Detail:  contentType,
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("form/legal-advise create rejected: invalid multipart body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "INVALID_MULTIPART_BODY",
			Message: "invalid multipart/form-data body",
			Detail:  err.Error(),
		})
	}

	actorUserID := ""
	timelineSource := service.TimelineSourcePublicForm
	if actor, ok := authActor(c); ok {
		actorUserID = actor.ID.String()
		timelineSource = service.TimelineSourceInternal
	}

	result, err := h.service.Create(c.Context(), service.CreateAsesoriaInput{
		AcceptsDataProcessing: formValue(form, "acceptsDataProcessing"),
		ConsultationDate:      formValue(form, "consultationDate"),
		FullName:              formValue(form, "fullName"),
		DocumentType:          formValue(form, "documentType"),
		OtherDocumentType:     formValue(form, "otherDocumentType"),
		DocumentNumber:        formValue(form, "documentNumber"),
		BirthDate:             formValue(form, "birthDate"),
		Age:                   formValue(form, "age"),
		MaritalStatus:         formValue(form, "maritalStatus"),
		OtherMaritalStatus:    formValue(form, "otherMaritalStatus"),
		Gender:                formValue(form, "gender"),
		Address:               formValue(form, "address"),
		HousingType:           formValue(form, "housingType"),
		OtherHousingType:      formValue(form, "otherHousingType"),
		SocioEconomicStratum:  formValue(form, "socioEconomicStratum"),
		SisbenCategory:        formValue(form, "sisbenCategory"),
		MobilePhone:           formValue(form, "mobilePhone"),
		Email:                 formValue(form, "email"),
		PopulationType:        formValue(form, "populationType"),
		OtherPopulationType:   formValue(form, "otherPopulationType"),
		HeadOfHousehold:       formValue(form, "headOfHousehold"),
		Occupation:            formValue(form, "occupation"),
		EducationLevel:        formValue(form, "educationLevel"),
		OtherEducationLevel:   formValue(form, "otherEducationLevel"),
		CaseDescription:       formValue(form, "caseDescription"),
		AuthorizesNotification: formValue(form,
			"authorizesNotification"),
		SubmittedAt:      formValue(form, "submittedAt"),
		IdentityDocument: formFile(form, "identityDocument"),
		UtilityBill:      formFile(form, "utilityBill"),
		ActorUserID:      actorUserID,
		TimelineSource:   timelineSource,
	})
	if err != nil {
		var validationErr *service.ValidationError
		if errors.As(err, &validationErr) {
			log.Printf("form/legal-advise validation failed: %+v", validationErr.Fields)
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "VALIDATION_ERROR",
				Message: validationErr.Error(),
				Fields:  validationErr.Fields,
			})
		}

		log.Printf("form/legal-advise create failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "LEGAL_ADVISE_CREATE_FAILED",
			Message: "failed to register legal advise form intake",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CreateAsesoriaResponse{
		ID:             result.ID,
		PreturnoNumber: result.PreturnoNumber,
		Status:         result.Status,
		CreatedAt:      result.CreatedAt,
	})
}

func (h *Handler) List(c *fiber.Ctx) error {
	page := parsePositiveInt(c.Query("page"), 1)
	limit := parsePositiveInt(c.Query("limit"), 20)

	result, err := h.service.List(c.Context(), service.ListPreturnosInput{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		log.Printf("form/legal-advise list failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "LEGAL_ADVISE_LIST_FAILED",
			Message: "failed to list preturnos",
			Detail:  err.Error(),
		})
	}

	return c.JSON(buildListPreturnosResponse(result))
}

func (h *Handler) ListAssigned(c *fiber.Ctx) error {
	actor, ok := authActor(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "missing authenticated session user",
		})
	}

	page := parsePositiveInt(c.Query("page"), 1)
	limit := parsePositiveInt(c.Query("limit"), 20)

	result, err := h.service.ListAssigned(c.Context(), service.ListAssignedPreturnosInput{
		Page:        page,
		Limit:       limit,
		ActorUserID: actor.ID.String(),
		ActorRoles:  actor.Roles,
	})
	if err != nil {
		if errors.Is(err, service.ErrAssignmentForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error:   "ASSIGNED_PRETURNOS_FORBIDDEN",
				Message: "user does not have permission to view assigned preturnos",
			})
		}

		log.Printf("form/legal-advise assigned list failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "LEGAL_ADVISE_ASSIGNED_LIST_FAILED",
			Message: "failed to list assigned preturnos",
			Detail:  err.Error(),
		})
	}

	return c.JSON(buildListPreturnosResponse(result))
}

func buildListPreturnosResponse(result service.ListPreturnosResult) dto.ListPreturnosResponse {
	items := make([]dto.PreturnoItem, 0, len(result.Items))
	for _, item := range result.Items {
		soportes := make([]dto.PreturnoSupport, 0, len(item.Soportes))
		for _, soporte := range item.Soportes {
			soportes = append(soportes, dto.PreturnoSupport{
				ID:          soporte.ID,
				Nombre:      soporte.Nombre,
				URL:         soporte.URL,
				ViewURL:     soporte.ViewURL,
				DownloadURL: soporte.DownloadURL,
				Tipo:        soporte.Tipo,
				Campo:       soporte.Campo,
				Tamano:      soporte.Tamano,
			})
		}

		timeline := make([]dto.PreturnoTimelineItem, 0, len(item.Timeline))
		for _, event := range item.Timeline {
			timeline = append(timeline, dto.PreturnoTimelineItem{
				ID:      event.ID,
				Fecha:   event.Fecha,
				Titulo:  event.Titulo,
				Detalle: event.Detalle,
				Usuario: event.Usuario,
				Source:  event.Source,
			})
		}

		items = append(items, dto.PreturnoItem{
			ID:                   item.ID,
			Radicado:             item.Radicado,
			Turno:                item.Turno,
			Estado:               item.Estado,
			AutorizacionDatos:    item.AutorizacionDatos,
			FechaConsulta:        item.FechaConsulta,
			NombreCompleto:       item.NombreCompleto,
			TipoDocumento:        item.TipoDocumento,
			NumeroDocumento:      item.NumeroDocumento,
			FechaNacimiento:      item.FechaNacimiento,
			Edad:                 item.Edad,
			EstadoCivil:          item.EstadoCivil,
			Genero:               item.Genero,
			Direccion:            item.Direccion,
			TipoVivienda:         item.TipoVivienda,
			Estrato:              item.Estrato,
			SisbenCategoria:      item.SisbenCategoria,
			Telefono:             item.Telefono,
			CorreoElectronico:    item.CorreoElectronico,
			TipoPoblacion:        item.TipoPoblacion,
			CabezaHogar:          item.CabezaHogar,
			Ocupacion:            item.Ocupacion,
			NivelEstudio:         item.NivelEstudio,
			Relato:               item.Relato,
			Soportes:             soportes,
			AutorizaNotificacion: item.AutorizaNotificacion,
			Timeline:             timeline,
			Citizen:              item.Citizen,
			Cedula:               item.Cedula,
		})
	}

	return dto.ListPreturnosResponse{
		Items: items,
		Total: result.Total,
		Page:  result.Page,
		Limit: result.Limit,
	}
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

func formFile(form *multipart.Form, key string) *multipart.FileHeader {
	if form == nil {
		return nil
	}

	if files, ok := form.File[key]; ok && len(files) > 0 {
		return files[0]
	}

	if files, ok := form.File[key+"[]"]; ok && len(files) > 0 {
		return files[0]
	}

	for fieldName, files := range form.File {
		if strings.HasPrefix(fieldName, key+"[") && len(files) > 0 {
			return files[0]
		}
	}

	return nil
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
