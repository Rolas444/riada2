package ports

import "github.com/riada2/internal/core/domain"

type AddressService interface {
	CreateOrUpdateAddressForUser(address *domain.Address, userID uint) (*domain.Address, error)
	DeleteAddressForUser(addressID uint, userID uint) error
}
