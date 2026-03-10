package service

import (
	"strings"
	"time"
)

func mapPublicStatus(dbStatus string) string {
	if strings.EqualFold(strings.TrimSpace(dbStatus), "RADICADA") {
		return statusPublicRecibido
	}
	return strings.TrimSpace(dbStatus)
}

func mapPublicListStatus(dbStatus string) string {
	switch strings.ToUpper(strings.TrimSpace(dbStatus)) {
	case "RADICADA":
		return "abierta"
	case "EN_GESTION":
		return "en_gestion"
	case "RESPONDIDA":
		return "respondida"
	case "CERRADA":
		return "cerrada"
	default:
		return strings.ToLower(strings.TrimSpace(dbStatus))
	}
}

func normalizeStatusFilter(value string) []string {
	normalized := normalizeText(value)
	switch normalized {
	case "":
		return nil
	case "abierta":
		return []string{"RADICADA", "EN_GESTION"}
	case "radicada":
		return []string{"RADICADA"}
	case "en_gestion", "engestion":
		return []string{"EN_GESTION"}
	case "respondida":
		return []string{"RESPONDIDA"}
	case "cerrada":
		return []string{"CERRADA"}
	default:
		return []string{strings.ToUpper(strings.TrimSpace(value))}
	}
}

func computeDeadline(fechaIngreso time.Time) time.Time {
	base := fechaIngreso.UTC()
	return time.Date(
		base.Year(),
		base.Month(),
		base.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	).AddDate(0, 0, 15)
}
