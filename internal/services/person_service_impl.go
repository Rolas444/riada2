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
	// Busca si ya existe un registro de persona para este UserID.
	existingPerson, err := s.personRepo.FindByID(*&person.ID)
	if err != nil {
		// Si el error es que no se encontró el registro, creamos uno nuevo.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := s.personRepo.Save(person)
			return person, err
		}
		// Si es otro tipo de error, lo devolvemos.
		return nil, err
	}

	// Si ya existe, actualizamos los campos del registro existente y lo guardamos.
	existingPerson.Name = person.Name
	existingPerson.MiddleName = person.MiddleName
	existingPerson.LastName = person.LastName
	existingPerson.Sex = person.Sex
	existingPerson.Birthday = person.Birthday
	existingPerson.DocNumber = person.DocNumber
	existingPerson.TypeDoc = person.TypeDoc

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
