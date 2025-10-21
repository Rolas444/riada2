package services

import (
	"errors"

	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
)

type MinistryMemberServiceImpl struct {
	repo           ports.MinistryMemberRepository
	membershipRepo ports.MembershipRepository
}

func NewMinistryMemberService(repo ports.MinistryMemberRepository, membershipRepo ports.MembershipRepository) ports.MinistryMemberService {
	return &MinistryMemberServiceImpl{repo: repo, membershipRepo: membershipRepo}
}

func (s *MinistryMemberServiceImpl) GetAll() ([]domain.MinistryMember, error) {
	return s.repo.FindAll()
}

func (s *MinistryMemberServiceImpl) Create(member *domain.MinistryMember) (*domain.MinistryMember, error) {
	// Validate membership exists for the person
	m, err := s.membershipRepo.FindByPersonID(member.PersonID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("no se puede agregar miembro sin membresía")
	}
	if member.Status == "" {
		member.Status = "A"
	}
	if err := s.repo.Save(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *MinistryMemberServiceImpl) Update(id uint, input *domain.MinistryMember) (*domain.MinistryMember, error) {
	// Not using id field in domain; pattern consistent with other services using repo lookups
	existing, err := s.repo.FindByMinistryAndPerson(input.MinistryID, input.PersonID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("ministry member not found")
	}
	existing.Role = input.Role
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
