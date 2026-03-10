package dto

type OverviewResponse struct {
	Stats          StatsSection         `json:"stats"`
	PendingQueue   []PendingQueueItem   `json:"pending_queue"`
	DeadlineAlerts []DeadlineAlertItem  `json:"deadline_alerts"`
	RecentActivity []RecentActivityItem `json:"recent_activity"`
}

type StatsSection struct {
	TurnosAtendidosHoy StatCard `json:"turnos_atendidos_hoy"`
	TurnosPendientes   StatCard `json:"turnos_pendientes"`
	CasosCreadosHoy    StatCard `json:"casos_creados_hoy"`
	AtencionesEnCurso  StatCard `json:"atenciones_en_curso"`
}

type StatCard struct {
	Value      int64  `json:"value"`
	Change     string `json:"change"`
	ChangeType string `json:"change_type"`
}

type PendingQueueItem struct {
	ID       string `json:"id"`
	Citizen  string `json:"citizen"`
	Motivo   string `json:"motivo"`
	Canal    string `json:"canal"`
	Time     string `json:"time"`
	Priority string `json:"priority"`
}

type DeadlineAlertItem struct {
	ID          string `json:"id"`
	Caso        string `json:"caso"`
	Description string `json:"description"`
	Responsible string `json:"responsible"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}

type RecentActivityItem struct {
	ID           string `json:"id"`
	User         string `json:"user"`
	Initials     string `json:"initials"`
	Action       string `json:"action"`
	Target       string `json:"target"`
	Time         string `json:"time"`
	Badge        string `json:"badge"`
	BadgeVariant string `json:"badge_variant"`
}

type OverviewQuery struct {
	Timezone       string
	StatsDate      string
	PendingLimit   int
	DeadlinesLimit int
	ActivityLimit  int
}
