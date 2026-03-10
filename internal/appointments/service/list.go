package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/appointments/repository"
	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, in ListPreturnosInput) (ListPreturnosResult, error) {
	if s.repo == nil {
		return ListPreturnosResult{}, errors.New("legal-advise repository is not initialized")
	}

	page := in.Page
	if page < 1 {
		page = 1
	}

	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	repoResult, err := s.repo.List(ctx, repository.ListInput{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return ListPreturnosResult{}, err
	}

	return mapPreturnosListResult(repoResult), nil
}

func (s *Service) ListAssigned(ctx context.Context, in ListAssignedPreturnosInput) (ListPreturnosResult, error) {
	if s.repo == nil {
		return ListPreturnosResult{}, errors.New("legal-advise repository is not initialized")
	}

	page := in.Page
	if page < 1 {
		page = 1
	}

	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	normalizedRoles := normalizeRoleCodes(in.ActorRoles)
	var assignedCoordinatorID *uuid.UUID

	switch {
	case canListAllAssignedPreturnos(normalizedRoles):
		assignedCoordinatorID = nil
	case containsRole(normalizedRoles, "COORDINADOR"):
		actorID := strings.TrimSpace(in.ActorUserID)
		if actorID == "" {
			return ListPreturnosResult{}, ErrAssignmentForbidden
		}
		parsedActorID, err := uuid.Parse(actorID)
		if err != nil {
			return ListPreturnosResult{}, ErrAssignmentForbidden
		}
		assignedCoordinatorID = &parsedActorID
	default:
		return ListPreturnosResult{}, ErrAssignmentForbidden
	}

	repoResult, err := s.repo.ListAssigned(ctx, repository.ListAssignedInput{
		Page:                  page,
		Limit:                 limit,
		AssignedCoordinatorID: assignedCoordinatorID,
	})
	if err != nil {
		return ListPreturnosResult{}, err
	}

	return mapPreturnosListResult(repoResult), nil
}

func mapPreturnosListResult(repoResult repository.ListResult) ListPreturnosResult {
	items := make([]PreturnoListItem, 0, len(repoResult.Items))
	for _, row := range repoResult.Items {
		payload := row.Payload
		nombreCompleto := strings.TrimSpace(payload.FullName)
		numeroDocumento := strings.TrimSpace(payload.DocumentNumber)

		soportes := make([]PreturnoSupport, 0, len(row.Attachments))
		for _, support := range row.Attachments {
			viewURL, downloadURL := supportLinks(support.FileObjectID, support.URL)
			displayURL := viewURL
			if displayURL == "" {
				displayURL = downloadURL
			}
			if displayURL == "" {
				displayURL = strings.TrimSpace(support.URL)
			}

			soportes = append(soportes, PreturnoSupport{
				ID:          support.ID,
				Nombre:      support.OriginalName,
				URL:         displayURL,
				ViewURL:     viewURL,
				DownloadURL: downloadURL,
				Tipo:        support.MimeType,
				Campo:       support.FieldName,
				Tamano:      support.SizeBytes,
			})
		}

		timeline := make([]PreturnoTimelineItem, 0, len(row.Timeline))
		for _, event := range row.Timeline {
			source := normalizeTimelineSource(event.Source)
			timeline = append(timeline, PreturnoTimelineItem{
				ID:      event.ID,
				Fecha:   event.CreatedAt.Format(time.RFC3339Nano),
				Titulo:  strings.TrimSpace(event.Title),
				Detalle: timelineDetailBySource(event.Title, event.Detail, source),
				Usuario: timelineEventUser(event.User, source),
				Source:  source,
			})
		}

		items = append(items, PreturnoListItem{
			ID:                   row.ID,
			Radicado:             row.PreturnoNumber,
			Turno:                "",
			Estado:               mapPreturnoStatusForFrontend(row.Status),
			AutorizacionDatos:    parseYesNoLoose(payload.AcceptsDataProcessing),
			FechaConsulta:        formatDateOnly(payload.ConsultationDate, row.ConsultationDate),
			NombreCompleto:       nombreCompleto,
			TipoDocumento:        selectWithOther(payload.DocumentType, payload.OtherDocumentType),
			NumeroDocumento:      numeroDocumento,
			FechaNacimiento:      formatDateOnly(payload.BirthDate, time.Time{}),
			Edad:                 parseIntLoose(payload.Age),
			EstadoCivil:          selectWithOther(payload.MaritalStatus, payload.OtherMaritalStatus),
			Genero:               strings.TrimSpace(payload.Gender),
			Direccion:            strings.TrimSpace(payload.Address),
			TipoVivienda:         selectWithOther(payload.HousingType, payload.OtherHousingType),
			Estrato:              parseOptionalIntLoose(payload.SocioEconomicStratum),
			SisbenCategoria:      strings.TrimSpace(payload.SisbenCategory),
			Telefono:             strings.TrimSpace(payload.MobilePhone),
			CorreoElectronico:    strings.TrimSpace(payload.Email),
			TipoPoblacion:        selectWithOther(payload.PopulationType, payload.OtherPopulationType),
			CabezaHogar:          parseOptionalYesNoLoose(payload.HeadOfHousehold),
			Ocupacion:            strings.TrimSpace(payload.Occupation),
			NivelEstudio:         selectWithOther(payload.EducationLevel, payload.OtherEducationLevel),
			Relato:               strings.TrimSpace(payload.CaseDescription),
			Soportes:             soportes,
			AutorizaNotificacion: parseYesNoLoose(payload.AuthorizesNotification),
			Timeline:             timeline,
			Citizen:              nombreCompleto,
			Cedula:               numeroDocumento,
		})
	}

	return ListPreturnosResult{
		Items: items,
		Total: repoResult.Total,
		Page:  repoResult.Page,
		Limit: repoResult.Limit,
	}
}

func canListAllAssignedPreturnos(roles []string) bool {
	return containsRole(roles, "SUPER_ADMIN") ||
		containsRole(roles, "ADMIN_CONSULTORIO") ||
		containsRole(roles, "SECRETARIA") ||
		containsRole(roles, "JEFATURA")
}

func containsRole(roles []string, role string) bool {
	target := strings.ToUpper(strings.TrimSpace(role))
	if target == "" {
		return false
	}

	for _, current := range roles {
		if strings.ToUpper(strings.TrimSpace(current)) == target {
			return true
		}
	}

	return false
}

func normalizeRoleCodes(roles []string) []string {
	seen := make(map[string]struct{}, len(roles))
	out := make([]string, 0, len(roles))

	for _, role := range roles {
		normalized := strings.ToUpper(strings.TrimSpace(role))
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}

	return out
}

func mapPreturnoStatusForFrontend(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case repository.StatusPendientePreturno:
		return "pendiente_clasificacion"
	case "PENDIENTE_PRETURNO":
		return "pendiente_clasificacion"
	case "EN_PROCESO":
		return "en_proceso"
	case "PROCESADO":
		return "procesado"
	case "RECHAZADO":
		return "rechazado"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func parseOptionalYesNoLoose(value string) *bool {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	normalized := normalizeText(value)
	switch normalized {
	case "yes", "si", "s", "true", "1":
		v := true
		return &v
	case "no", "n", "false", "0":
		v := false
		return &v
	default:
		return nil
	}
}

func parseOptionalIntLoose(value string) *int {
	out, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return nil
	}
	return &out
}

func parseYesNoLoose(value string) bool {
	parsed := parseOptionalYesNoLoose(value)
	return parsed != nil && *parsed
}

func parseIntLoose(value string) int {
	parsed := parseOptionalIntLoose(value)
	if parsed == nil {
		return 0
	}
	return *parsed
}

func normalizeTimelineSource(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case repository.TimelineSourcePublicForm:
		return repository.TimelineSourcePublicForm
	case repository.TimelineSourceInternal:
		return repository.TimelineSourceInternal
	default:
		return repository.TimelineSourceInternal
	}
}

func timelineEventUser(user string, source string) string {
	name := strings.TrimSpace(user)
	if name != "" {
		return name
	}
	if source == repository.TimelineSourcePublicForm {
		return "Formulario publico"
	}
	return "Sistema"
}

func timelineDetailBySource(title string, detail string, source string) string {
	titleText := strings.TrimSpace(title)
	text := strings.TrimSpace(detail)

	if strings.EqualFold(titleText, "Pre-turno registrado") {
		switch source {
		case repository.TimelineSourcePublicForm:
			if text == "" {
				return "Registro automatico desde formulario"
			}
			return text
		case repository.TimelineSourceInternal:
			if text == "" || strings.EqualFold(text, "Registro automatico desde formulario") {
				return "Registro manual desde preturnos"
			}
			return text
		default:
			return text
		}
	}

	return text
}

func selectWithOther(value string, other string) string {
	main := strings.TrimSpace(value)
	alt := strings.TrimSpace(other)
	normalized := normalizeText(main)
	if normalized == "otro" || normalized == "otra" || normalized == "other" {
		if alt != "" {
			return alt
		}
	}
	if main != "" {
		return main
	}
	return alt
}

func formatDateOnly(raw string, fallback time.Time) string {
	value := strings.TrimSpace(raw)
	if value == "" && !fallback.IsZero() {
		return fallback.UTC().Format("2006-01-02")
	}
	if value == "" {
		return ""
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			if layout == "2006-01-02" {
				return parsed.Format("2006-01-02")
			}
			return parsed.UTC().Format("2006-01-02")
		}
	}

	return value
}

func supportLinks(fileObjectID string, fallbackURL string) (string, string) {
	id := strings.TrimSpace(fileObjectID)
	if id != "" {
		base := "/api/v1/documents/files/" + id
		return base + "/view", base + "/download"
	}

	url := strings.TrimSpace(fallbackURL)
	if url == "" {
		return "", ""
	}

	if strings.HasSuffix(url, "/download") {
		return strings.TrimSuffix(url, "/download") + "/view", url
	}
	if strings.HasSuffix(url, "/view") {
		return url, strings.TrimSuffix(url, "/view") + "/download"
	}

	return url, url
}
