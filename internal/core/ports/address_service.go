package ports

import "github.com/riada2/internal/core/domain"

type AddressService interface {
	CreateOrUpdateAddress(address *domain.Address) (*domain.Address, error)
	DeleteAddress(addressID uint) error
}
