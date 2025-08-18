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

func (s *addressServiceImpl) CreateOrUpdateAddress(address *domain.Address) (*domain.Address, error) {

	// Validar que la persona (PersonID) a la que se asocia la dirección realmente exista.
	_, err := s.personRepo.FindByID(address.PersonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Usamos un error específico para que el handler pueda interpretarlo.
			return nil, ports.ErrPersonNotFound
		}
		return nil, err // Otro error de base de datos.
	}

	// Es una dirección nueva, verificar el límite. Una persona no puede tener más de 3 direcciones.
	count, err := s.addressRepo.CountByPersonID(address.PersonID)
	if err != nil {
		return nil, err
	}
	if count >= 2 {
		return nil, errors.New("una persona no puede tener más de 2 direcciones")
	}

	if err := s.addressRepo.Save(address); err != nil {
		return nil, err
	}

	return address, nil
}

func (s *addressServiceImpl) DeleteAddress(addressID uint) error {
	// Find the person associated with the user

	return s.addressRepo.Delete(addressID)
}
