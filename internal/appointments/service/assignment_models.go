package service

import "time"

type AssignPreturnoInput struct {
	PreturnoID    string
	CoordinatorID string
	ServiceTypeID string
	Observations  string
	ActorUserID   string
	ActorRoles    []string
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
	ServiceTypeID         string
	TimelineEvent         AssignmentTimelineEvent
}

type AssignmentCoordinatorOption struct {
	ID          string
	DisplayName string
	Email       string
}

type AssignmentServiceTypeOption struct {
	ID   string
	Code string
	Name string
}

type AssignmentOptionsResult struct {
	Coordinators []AssignmentCoordinatorOption
	ServiceTypes []AssignmentServiceTypeOption
}
