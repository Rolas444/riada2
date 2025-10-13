package ports

import "github.com/riada2/internal/core/domain"

type MinistryService interface {
	Create(ministry *domain.Ministry) (*domain.Ministry, error)
	Update(id uint, ministry *domain.Ministry) (*domain.Ministry, error)
}
