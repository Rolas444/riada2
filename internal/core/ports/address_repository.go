package ports

import "github.com/riada2/internal/core/domain"

type AddressRepository interface {
	Save(address *domain.Address) error
	Delete(id uint) error
	FindByID(id uint) (*domain.Address, error)
	CountByPersonID(personID uint) (int64, error)
}
