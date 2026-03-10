package repository

import "errors"

var (
	ErrPreturnoNotFound           = errors.New("preturno not found")
	ErrCoordinatorNotFound        = errors.New("coordinator not found")
	ErrServiceTypeNotFound        = errors.New("service type not found")
	ErrPreturnoAssignmentConflict = errors.New("preturno does not allow assignment")
)
