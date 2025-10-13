package ports

import "github.com/riada2/internal/core/domain"

type MinistryMemberService interface {
	Create(member *domain.MinistryMember) (*domain.MinistryMember, error)
	Update(id uint, member *domain.MinistryMember) (*domain.MinistryMember, error)
}
