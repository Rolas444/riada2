package ports

import "github.com/riada2/internal/core/domain"

type MinistryMemberRepository interface {
	FindAll() ([]domain.MinistryMember, error)
    FindByMinistryID(ministryID uint) ([]domain.MinistryMember, error)
	Save(member *domain.MinistryMember) error
	Update(member *domain.MinistryMember) error
	FindByMinistryAndPerson(ministryID, personID uint) (*domain.MinistryMember, error)
}
