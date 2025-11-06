package ports

import "github.com/riada2/internal/core/domain"

type MinistryMemberService interface {
	GetAll() ([]domain.MinistryMember, error)
    GetByMinistryID(ministryID uint) ([]domain.MinistryMember, error)
	Create(member *domain.MinistryMember) (*domain.MinistryMember, error)
	Update(id uint, member *domain.MinistryMember) (*domain.MinistryMember, error)
}
