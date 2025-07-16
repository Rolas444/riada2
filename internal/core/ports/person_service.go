package ports

import "github.com/riada2/internal/core/domain"

type PersonService interface {
	CreateOrUpdatePersonForUser(person *domain.Person) (*domain.Person, error)
	CreatePerson(person *domain.Person) (*domain.Person, error)
	DeletePerson(id uint) error
	GetPersonByID(id uint) (*domain.Person, error)
	SearchPersons(searchTerm string) ([]domain.Person, error)
}
