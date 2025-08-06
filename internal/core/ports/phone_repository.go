package ports

import "github.com/riada2/internal/core/domain"

// PhoneRepository defines the methods for interacting with phone storage.
// This interface is the port for the phone repository adapter.
type PhoneRepository interface {
	FindByID(id uint) (*domain.Phone, error)
	Save(phone *domain.Phone) error
	Delete(id uint) error
	// CountByPersonID counts the number of phones associated with a person.
	CountByPersonID(personID uint) (int64, error)
}
