package ports

import "github.com/riada2/internal/core/domain"

// PersonRepository defines the methods that any
// data storage provider needs to implement to get and store persons.
type PersonRepository interface {
	FindByUserID(userID uint) (*domain.Person, error)
	Save(person *domain.Person) error
	Delete(id uint) error
	FindByID(id uint) (*domain.Person, error)
	Search(searchTerm string) ([]domain.Person, error)
}
