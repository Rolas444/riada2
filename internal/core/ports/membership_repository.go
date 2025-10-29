package ports

import "github.com/riada2/internal/core/domain"

// MembershipRepository es el puerto para la persistencia de membresías.
type MembershipRepository interface {
	Save(membership *domain.Membership) error
	FindByPersonID(personID uint) (*domain.Membership, error)
	FindByID(id uint) (*domain.Membership, error)
	FindAll() ([]domain.Membership, error)
	Update(membership *domain.Membership) error
	Delete(id uint) error
}
