package ports

import "github.com/riada2/internal/core/domain"

// MembershipService es el puerto para la lógica de negocio de membresías.
type MembershipService interface {
	CreateMembership(personID uint, membershipData *domain.Membership) (*domain.Membership, error)
	GetMembershipByPersonID(personID uint) (*domain.Membership, error)
	GetMembershipByID(id uint) (*domain.Membership, error)
	GetAllMemberships() ([]domain.Membership, error)
	UpdateMembership(id uint, membershipData *domain.Membership) (*domain.Membership, error)
	DeleteMembership(id uint) error
}
