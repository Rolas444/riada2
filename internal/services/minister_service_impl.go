package services

import (
	"errors"

	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
)

type MinistryServiceImpl struct {
	repo ports.MinistryRepository
}

func NewMinistryService(repo ports.MinistryRepository) ports.MinistryService {
	return &MinistryServiceImpl{repo: repo}
}

func (s *MinistryServiceImpl) GetAll() ([]domain.Ministry, error) {
	return s.repo.FindAll()
}

func (s *MinistryServiceImpl) Create(ministry *domain.Ministry) (*domain.Ministry, error) {
	if ministry.Name == "" {
		return nil, errors.New("name is required")
	}
	if ministry.Status == "" {
		ministry.Status = "A"
	}
	if err := s.repo.Save(ministry); err != nil {
		return nil, err
	}
	return ministry, nil
}

func (s *MinistryServiceImpl) Update(id uint, input *domain.Ministry) (*domain.Ministry, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("ministry not found")
	}
	existing.Name = input.Name
	existing.Description = input.Description
	existing.Mission = input.Mission
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
