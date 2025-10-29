package ports

import "github.com/riada2/internal/core/domain"

type MinistryRepository interface {
	FindAll() ([]domain.Ministry, error)
	Save(ministry *domain.Ministry) error
	Update(ministry *domain.Ministry) error
	FindByID(id uint) (*domain.Ministry, error)
}
