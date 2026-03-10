package dto

import "time"

type AssignPreturnoRequest struct {
	CoordinatorID string `json:"coordinator_id"`
	ServiceTypeID string `json:"service_type_id"`
	Observations  string `json:"observations"`
}

type AssignmentTimelineEvent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

type AssignPreturnoResponse struct {
	ID                    string                  `json:"id"`
	Status                string                  `json:"status"`
	AssignedCoordinatorID string                  `json:"assigned_coordinator_id"`
	ServiceTypeID         string                  `json:"service_type_id"`
	TimelineEvent         AssignmentTimelineEvent `json:"timeline_event"`
}

type AssignmentCoordinatorOption struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

type AssignmentServiceTypeOption struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type AssignmentOptionsResponse struct {
	Coordinators []AssignmentCoordinatorOption `json:"coordinators"`
	ServiceTypes []AssignmentServiceTypeOption `json:"service_types"`
}
