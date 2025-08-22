package services

import (
	"errors"

	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
)

type MembershipServiceImpl struct {
	membershipRepo ports.MembershipRepository
}

func NewMembershipService(membershipRepo ports.MembershipRepository) ports.MembershipService {
	return &MembershipServiceImpl{
		membershipRepo: membershipRepo,
	}
}

func (s *MembershipServiceImpl) CreateMembership(personID uint, membershipData *domain.Membership) (*domain.Membership, error) {
	// Verificar que la persona no tenga ya una membresía
	existingMembership, err := s.membershipRepo.FindByPersonID(personID)
	if err != nil {
		return nil, err
	}
	if existingMembership != nil {
		return nil, errors.New("person already has a membership")
	}

	// Establecer valores por defecto si no se proporcionan
	if membershipData.State == "" {
		membershipData.State = domain.Active
	}

	// Guardar la membresía
	err = s.membershipRepo.Save(membershipData)
	if err != nil {
		return nil, err
	}

	return membershipData, nil
}

func (s *MembershipServiceImpl) GetMembershipByPersonID(personID uint) (*domain.Membership, error) {
	return s.membershipRepo.FindByPersonID(personID)
}

func (s *MembershipServiceImpl) GetMembershipByID(id uint) (*domain.Membership, error) {
	return s.membershipRepo.FindByID(id)
}

func (s *MembershipServiceImpl) GetAllMemberships() ([]domain.Membership, error) {
	return s.membershipRepo.FindAll()
}

func (s *MembershipServiceImpl) UpdateMembership(id uint, membershipData *domain.Membership) (*domain.Membership, error) {
	// Verificar que la membresía existe
	existingMembership, err := s.membershipRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existingMembership == nil {
		return nil, errors.New("membership not found")
	}

	// Actualizar los campos
	existingMembership.StartedAt = membershipData.StartedAt
	existingMembership.MembershipSigned = membershipData.MembershipSigned
	existingMembership.State = membershipData.State
	existingMembership.Transferred = membershipData.Transferred
	existingMembership.NameLastChurch = membershipData.NameLastChurch
	existingMembership.Baptized = membershipData.Baptized
	existingMembership.BaptismDate = membershipData.BaptismDate

	// Guardar los cambios
	err = s.membershipRepo.Update(existingMembership)
	if err != nil {
		return nil, err
	}

	return existingMembership, nil
}

func (s *MembershipServiceImpl) DeleteMembership(id uint) error {
	// Verificar que la membresía existe
	existingMembership, err := s.membershipRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingMembership == nil {
		return errors.New("membership not found")
	}

	return s.membershipRepo.Delete(id)
}
