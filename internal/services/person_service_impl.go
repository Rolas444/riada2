package services

import (
	"errors"

	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type personServiceImpl struct {
	personRepo ports.PersonRepository
}

func NewPersonService(personRepo ports.PersonRepository) ports.PersonService {
	return &personServiceImpl{personRepo}
}

// checkDocumentUniqueness valida que la combinación de TypeDoc y DocNumber sea única.
// NOTA: Esto asume que un nuevo método `FindByDocument` existe en la interfaz `ports.PersonRepository`.
func (s *personServiceImpl) checkDocumentUniqueness(person *domain.Person) error {
	if person.TypeDoc == nil || person.DocNumber == nil || *person.DocNumber == "" {
		return nil // No hay documento para verificar, se omite la validación.
	}

	existing, err := s.personRepo.FindByDocument(*person.TypeDoc, *person.DocNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // El documento no existe, lo cual es correcto.
		}
		return err // Otro error de base de datos.
	}

	// Si se encontró una persona y su ID es diferente al de la persona que se está guardando.
	// Esto cubre tanto la creación (person.ID es 0) como la actualización (person.ID != existing.ID).
	if existing.ID != person.ID {
		return ports.ErrPersonDocumentExists
	}

	return nil
}

func (s *personServiceImpl) CreateOrUpdatePersonForUser(person *domain.Person) (*domain.Person, error) {
	if person.UserID == nil {
		return nil, errors.New("UserID is required to create or update a person for a user")
	}

	// Si el ID de la persona es 0, es una creación.
	if person.ID == 0 {
		if err := s.checkDocumentUniqueness(person); err != nil {
			return nil, err
		}
		err := s.personRepo.Save(person)
		return person, err
	}

	// Si se proporciona un ID, es una operación de actualización.
	existingPerson, err := s.personRepo.FindByID(person.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrPersonNotFound
		}
		return nil, err
	}

	// Comprobación de seguridad:
	if existingPerson.UserID == nil || *existingPerson.UserID != *person.UserID {
		return nil, errors.New("authorization failed: you can only update your own person record")
	}

	// Actualizamos los campos del registro existente en memoria.
	existingPerson.Name = person.Name
	existingPerson.MiddleName = person.MiddleName
	existingPerson.LastName = person.LastName
	existingPerson.Sex = person.Sex
	existingPerson.Birthday = person.Birthday
	existingPerson.DocNumber = person.DocNumber
	existingPerson.TypeDoc = person.TypeDoc
	existingPerson.Email = person.Email
	existingPerson.Photo = person.Photo

	// Validar la unicidad del documento DESPUÉS de actualizar los campos y ANTES de guardar.
	if err := s.checkDocumentUniqueness(existingPerson); err != nil {
		return nil, err
	}

	// Guardamos la entidad actualizada.
	err = s.personRepo.Save(existingPerson)
	return existingPerson, err
}

func (s *personServiceImpl) CreatePerson(person *domain.Person) (*domain.Person, error) {
	if err := s.checkDocumentUniqueness(person); err != nil {
		return nil, err
	}
	err := s.personRepo.Save(person)
	return person, err
}

func (s *personServiceImpl) DeletePerson(id uint) error {
	return s.personRepo.Delete(id)
}

func (s *personServiceImpl) GetPersonByID(id uint) (*domain.Person, error) {
	return s.personRepo.FindByID(id)
}

func (s *personServiceImpl) SearchPersons(searchTerm string) ([]domain.Person, error) {
	return s.personRepo.Search(searchTerm)
}
