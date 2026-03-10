package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	TimelineSourcePublicForm = "formulario_publico"
	TimelineSourceInternal   = "lista_interna"
)

func parseOptionalDate(value string) *time.Time {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return nil
	}

	layouts := []string{
		"2006-01-02",
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, raw)
		if err == nil {
			normalized := parsed.UTC()
			return &normalized
		}
	}
	return nil
}

func parseOptionalSmallInt(value string) *int16 {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return nil
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}
	if parsed < 0 || parsed > 32767 {
		return nil
	}
	v := int16(parsed)
	return &v
}

func normalizeStratum(value string) *int16 {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return nil
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}
	if parsed < 1 || parsed > 6 {
		return nil
	}
	v := int16(parsed)
	return &v
}

func normalizeEmailOrPlaceholder(primary string, fallback string, seed string) string {
	candidate := strings.TrimSpace(primary)
	if candidate == "" {
		candidate = strings.TrimSpace(fallback)
	}
	if strings.Contains(candidate, "@") {
		return candidate
	}
	return "no-informado+" + shortID(seed) + "@sicou.local"
}

func shortID(value string) string {
	id := strings.TrimSpace(value)
	id = strings.ReplaceAll(id, "-", "")
	if len(id) >= 8 {
		return strings.ToLower(id[:8])
	}
	if id == "" {
		return "na"
	}
	return strings.ToLower(id)
}

func selectWithOther(value string, other string) string {
	mainValue := strings.TrimSpace(value)
	otherValue := strings.TrimSpace(other)
	if mainValue == "" {
		return otherValue
	}
	normalized := strings.ToLower(mainValue)
	switch normalized {
	case "otro", "otra", "other":
		if otherValue != "" {
			return otherValue
		}
	}
	return mainValue
}

func formatPreturnoNumber(createdAt time.Time, sequence int64) string {
	return fmt.Sprintf("PT-%d-%04d", createdAt.Year(), sequence)
}

func buildCitizenNotificationSubject(preturnoNumber string) string {
	number := strings.TrimSpace(preturnoNumber)
	if number == "" {
		number = "PT-0000-0000"
	}
	return fmt.Sprintf("Confirmacion de preturno - %s", number)
}

func buildCitizenNotificationBody(preturnoNumber string, citizenName string) string {
	number := strings.TrimSpace(preturnoNumber)
	if number == "" {
		number = "PT-0000-0000"
	}

	name := strings.TrimSpace(citizenName)
	if name == "" {
		name = "usuario"
	}

	return fmt.Sprintf(
		"Hola %s,\n\nGracias por registrar tu solicitud de asesoria juridica.\nTu numero de preturno es %s.\n\nCon este numero puedes hacer seguimiento a tu solicitud.\n\nAtentamente,\nConsultorio Juridico UNIMETA\n",
		name,
		number,
	)
}

func normalizeTimelineSource(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case TimelineSourcePublicForm:
		return TimelineSourcePublicForm
	case TimelineSourceInternal:
		return TimelineSourceInternal
	default:
		return TimelineSourceInternal
	}
}

func buildPreturnoCreatedDetail(source string) string {
	switch normalizeTimelineSource(source) {
	case TimelineSourcePublicForm:
		return "Registro automatico desde formulario"
	default:
		return "Registro manual desde preturnos"
	}
}
