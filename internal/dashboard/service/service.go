package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/dashboard/dto"
	"github.com/DuvanRozoParra/sicou/internal/dashboard/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Overview(ctx context.Context, day time.Time, loc *time.Location, pendingLimit, deadlinesLimit, activityLimit int) (dto.OverviewResponse, error) {
	response := dto.OverviewResponse{
		Stats: dto.StatsSection{
			TurnosAtendidosHoy: dto.StatCard{Value: 0, Change: "0%", ChangeType: "positive"},
			TurnosPendientes:   dto.StatCard{Value: 0, Change: "0%", ChangeType: "negative"},
			CasosCreadosHoy:    dto.StatCard{Value: 0, Change: "0%", ChangeType: "positive"},
			AtencionesEnCurso:  dto.StatCard{Value: 0, Change: "0%", ChangeType: "positive"},
		},
		PendingQueue:   make([]dto.PendingQueueItem, 0),
		DeadlineAlerts: make([]dto.DeadlineAlertItem, 0),
		RecentActivity: make([]dto.RecentActivityItem, 0),
	}

	stats, err := s.repo.OverviewStats(ctx, day, loc)
	if err == nil {
		response.Stats.TurnosAtendidosHoy.Value = stats.TurnosAtendidosHoyCurrent
		response.Stats.TurnosAtendidosHoy.Change = percentageChange(stats.TurnosAtendidosHoyCurrent, stats.TurnosAtendidosHoyPrev)

		response.Stats.TurnosPendientes.Value = stats.TurnosPendientesCurrent
		response.Stats.TurnosPendientes.Change = percentageChange(stats.TurnosPendientesCurrent, stats.TurnosPendientesPrev)

		response.Stats.CasosCreadosHoy.Value = stats.CasosCreadosHoyCurrent
		response.Stats.CasosCreadosHoy.Change = percentageChange(stats.CasosCreadosHoyCurrent, stats.CasosCreadosHoyPrev)

		response.Stats.AtencionesEnCurso.Value = stats.AtencionesEnCursoCurrent
		response.Stats.AtencionesEnCurso.Change = percentageChange(stats.AtencionesEnCursoCurrent, stats.AtencionesEnCursoPrev)
	} else {
		log.Printf("dashboard overview: stats query failed: %v", err)
	}

	pendingRows, err := s.repo.PendingQueue(ctx, pendingLimit)
	if err == nil {
		for _, row := range pendingRows {
			response.PendingQueue = append(response.PendingQueue, dto.PendingQueueItem{
				ID:       fmt.Sprintf("T-%04d", row.Consecutive),
				Citizen:  strings.TrimSpace(row.Citizen),
				Motivo:   compactText(row.Motivo, 110),
				Canal:    strings.TrimSpace(row.Canal),
				Time:     elapsedHHMM(row.CreatedAt.In(loc), time.Now().In(loc)),
				Priority: priorityLabel(row.Priority),
			})
		}
	} else {
		log.Printf("dashboard overview: pending queue query failed: %v", err)
	}

	deadlineRows, err := s.repo.DeadlineAlerts(ctx, day, deadlinesLimit)
	if err == nil {
		for i, row := range deadlineRows {
			response.DeadlineAlerts = append(response.DeadlineAlerts, dto.DeadlineAlertItem{
				ID:          fmt.Sprintf("CASE-%03d", i+1),
				Caso:        safeDefault(row.CaseCode, "EXP-0000-0000"),
				Description: "Vence termino de respuesta",
				Responsible: safeDefault(row.Responsible, "Sin responsable"),
				DueDate:     row.DueDate.Format("2006-01-02"),
				Status:      mapDeadlineStatus(row.TermStatus),
			})
		}
	} else {
		log.Printf("dashboard overview: deadline alerts query failed: %v", err)
	}

	activityRows, err := s.repo.RecentActivity(ctx, activityLimit)
	if err == nil {
		for i, row := range activityRows {
			action, variant := actionAndVariant(row.EventType)
			user := safeDefault(strings.TrimSpace(row.UserName), "Sistema")
			response.RecentActivity = append(response.RecentActivity, dto.RecentActivityItem{
				ID:           fmt.Sprintf("evt_%d", i+1),
				User:         user,
				Initials:     initials(user),
				Action:       action,
				Target:       fmt.Sprintf("caso %s", safeDefault(row.CaseCode, "EXP-0000-0000")),
				Time:         agoText(row.CreatedAt.In(loc), time.Now().In(loc)),
				Badge:        "Casos",
				BadgeVariant: variant,
			})
		}
	} else {
		log.Printf("dashboard overview: recent activity query failed: %v", err)
	}

	return response, nil
}

func percentageChange(current, previous int64) string {
	if previous == 0 {
		if current == 0 {
			return "0%"
		}
		return "100%"
	}

	change := (float64(current-previous) / float64(previous)) * 100
	value := int64(math.Round(change))
	return fmt.Sprintf("%d%%", value)
}

func elapsedHHMM(start, now time.Time) string {
	if now.Before(start) {
		return "00:00"
	}
	d := now.Sub(start)
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if hours < 0 {
		hours = 0
	}
	if minutes < 0 {
		minutes = 0
	}
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func priorityLabel(value int) string {
	switch {
	case value >= 2:
		return "alta"
	case value == 1:
		return "media"
	default:
		return "baja"
	}
}

func mapDeadlineStatus(in string) string {
	switch strings.ToUpper(strings.TrimSpace(in)) {
	case "VENCIDO":
		return "vencido"
	case "PROXIMO":
		return "por_vencer"
	default:
		return "en_tiempo"
	}
}

func actionAndVariant(eventType string) (string, string) {
	switch strings.ToUpper(strings.TrimSpace(eventType)) {
	case "CREATED", "CASE_CREATED":
		return "creo", "default"
	case "DOCUMENT_ADDED", "ATTACHMENT_ADDED":
		return "adjunto", "outline"
	case "ASSIGNMENT", "ASSIGNED":
		return "asigno", "secondary"
	default:
		return "actualizo", "secondary"
	}
}

func agoText(t, now time.Time) string {
	if now.Before(t) {
		return "hace 0 min"
	}
	d := now.Sub(t)
	if d < time.Minute {
		return "hace 0 min"
	}
	if d < time.Hour {
		return fmt.Sprintf("hace %d min", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("hace %d h", int(d.Hours()))
	}
	return fmt.Sprintf("hace %d d", int(d.Hours()/24))
}

func initials(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return "NA"
	}
	if len(parts) == 1 {
		r := []rune(parts[0])
		if len(r) == 0 {
			return "NA"
		}
		return strings.ToUpper(string(r[0]))
	}
	first := []rune(parts[0])
	last := []rune(parts[len(parts)-1])
	if len(first) == 0 || len(last) == 0 {
		return "NA"
	}
	return strings.ToUpper(string(first[0]) + string(last[0]))
}

func compactText(in string, max int) string {
	value := strings.Join(strings.Fields(strings.TrimSpace(in)), " ")
	if value == "" {
		return "Sin descripcion"
	}
	if max <= 0 || len(value) <= max {
		return value
	}
	if max <= 3 {
		return value[:max]
	}
	return value[:max-3] + "..."
}

func safeDefault(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}
