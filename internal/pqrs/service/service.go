package service

import "github.com/DuvanRozoParra/sicou/internal/pqrs/repository"

type Service struct {
	repo        repository.Repository
	fileObjects FileObjectService
}

func NewService(repo repository.Repository, fileObjects FileObjectService) *Service {
	return &Service{
		repo:        repo,
		fileObjects: fileObjects,
	}
}
