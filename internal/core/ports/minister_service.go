package ports

import "github.com/riada2/internal/core/domain"

type MinistryService interface {
	GetAll() ([]domain.Ministry, error)
	Create(ministry *domain.Ministry) (*domain.Ministry, error)
	Update(id uint, ministry *domain.Ministry) (*domain.Ministry, error)
}
