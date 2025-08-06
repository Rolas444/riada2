package services

import (
	"errors"

	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type addressServiceImpl struct {
	addressRepo ports.AddressRepository
	personRepo  ports.PersonRepository
}

func NewAddressService(addressRepo ports.AddressRepository, personRepo ports.PersonRepository) ports.AddressService {
	return &addressServiceImpl{addressRepo, personRepo}
}

func (s *addressServiceImpl) CreateOrUpdateAddressForUser(address *domain.Address, userID uint) (*domain.Address, error) {

	// Es una dirección nueva, verificar el límite. Una persona no puede tener más de 3 direcciones.
	count, err := s.addressRepo.CountByPersonID(address.PersonID)
	if err != nil {
		return nil, err
	}
	if count >= 3 {
		return nil, errors.New("una persona no puede tener más de 3 direcciones")
	}

	if err := s.addressRepo.Save(address); err != nil {
		return nil, err
	}

	return address, nil
}

func (s *addressServiceImpl) DeleteAddressForUser(addressID uint, userID uint) error {
	// Find the person associated with the user
	person, err := s.personRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("no person profile found for this user")
	}

	// Find the address to be deleted
	address, err := s.addressRepo.FindByID(addressID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("address not found")
		}
		return err
	}

	// Check if the address belongs to the correct person
	if address.PersonID != person.ID {
		return errors.New("authorization failed: you can only delete your own addresses")
	}

	return s.addressRepo.Delete(addressID)
}
