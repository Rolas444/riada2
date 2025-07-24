package ports

import (
	"errors"

	"github.com/riada2/internal/core/domain"
)

var ErrPersonNotFound = errors.New("person not found")

type PersonService interface {
	CreateOrUpdatePersonForUser(person *domain.Person) (*domain.Person, error)
	CreatePerson(person *domain.Person) (*domain.Person, error)
	DeletePerson(id uint) error
	GetPersonByID(id uint) (*domain.Person, error)
	SearchPersons(searchTerm string) ([]domain.Person, error)
}
