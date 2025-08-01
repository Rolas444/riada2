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

func (s *personServiceImpl) CreateOrUpdatePersonForUser(person *domain.Person) (*domain.Person, error) {
	if person.UserID == nil {
		return nil, errors.New("UserID is required to create or update a person for a user")
	}

	// Si el ID de la persona es 0, se asume que es una operación de creación.
	if person.ID == 0 {
		// GORM creará un nuevo registro porque la clave primaria (ID) es su valor cero.
		err := s.personRepo.Save(person)
		return person, err
	}

	// Si se proporciona un ID, es una operación de actualización.
	// Primero, buscamos el registro existente en la base de datos.
	existingPerson, err := s.personRepo.FindByID(person.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// El cliente intentó actualizar una persona que no existe.
			return nil, ports.ErrPersonNotFound
		}
		// Otro error de base de datos.
		return nil, err
	}

	// Comprobación de seguridad: Asegurarse de que el usuario autenticado
	// solo pueda modificar su propio registro de persona.
	if existingPerson.UserID == nil || *existingPerson.UserID != *person.UserID {
		return nil, errors.New("authorization failed: you can only update your own person record")
	}

	// Actualizamos los campos del registro existente con los nuevos datos.
	existingPerson.Name = person.Name
	existingPerson.MiddleName = person.MiddleName
	existingPerson.LastName = person.LastName
	existingPerson.Sex = person.Sex
	existingPerson.Birthday = person.Birthday
	existingPerson.DocNumber = person.DocNumber
	existingPerson.TypeDoc = person.TypeDoc
	existingPerson.Email = person.Email
	existingPerson.Photo = person.Photo
	existingPerson.UserID = person.UserID // Actualiza quién hizo la última modificación

	// Guardamos la entidad actualizada. GORM realizará un UPDATE.
	err = s.personRepo.Save(existingPerson)
	return existingPerson, err
}

func (s *personServiceImpl) CreatePerson(person *domain.Person) (*domain.Person, error) {
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
