package services

import (
	"errors"

	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type phoneServiceImpl struct {
	phoneRepo  ports.PhoneRepository
	personRepo ports.PersonRepository
}

func NewPhoneService(phoneRepo ports.PhoneRepository, personRepo ports.PersonRepository) ports.PhoneService {
	return &phoneServiceImpl{phoneRepo, personRepo}
}

func (s *phoneServiceImpl) CreateOrUpdatePhoneForUser(phone *domain.Phone, userID uint) (*domain.Phone, error) {
	// Find the person associated with the user
	person, err := s.personRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no person profile found for this user, please create one first")
		}
		return nil, err
	}

	// If it's an update, check ownership
	if phone.ID != 0 {
		existingPhone, err := s.phoneRepo.FindByID(phone.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("phone not found")
			}
			return nil, err
		}
		// Check if the phone belongs to the correct person
		if existingPhone.PersonID != person.ID {
			return nil, errors.New("authorization failed: you can only update your own phones")
		}
	} else {
		// Es un teléfono nuevo, verificar el límite. Una persona no puede tener más de 2 números.
		// Se asume que el repositorio provee un método para contar los teléfonos de una persona.
		count, err := s.phoneRepo.CountByPersonID(person.ID)
		if err != nil {
			return nil, err
		}
		if count >= 2 {
			return nil, errors.New("una persona no puede tener más de 2 números de teléfono")
		}
	}

	// Associate the phone with the person and save
	phone.PersonID = person.ID
	if err := s.phoneRepo.Save(phone); err != nil {
		return nil, err
	}

	return phone, nil
}

func (s *phoneServiceImpl) DeletePhoneForUser(phoneID uint, userID uint) error {
	// Find the person associated with the user
	person, err := s.personRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("no person profile found for this user")
	}

	// Find the phone to be deleted
	phone, err := s.phoneRepo.FindByID(phoneID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("phone not found")
		}
		return err
	}

	// Check if the phone belongs to the correct person
	if phone.PersonID != person.ID {
		return errors.New("authorization failed: you can only delete your own phones")
	}

	return s.phoneRepo.Delete(phoneID)
}
