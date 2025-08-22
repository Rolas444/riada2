package repository

import (
	"errors"

	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type GormMembershipRepository struct {
	db *gorm.DB
}

func NewGormMembershipRepository(db *gorm.DB) ports.MembershipRepository {
	return &GormMembershipRepository{db: db}
}

func (r *GormMembershipRepository) Save(membership *domain.Membership) error {
	return r.db.Create(membership).Error
}

func (r *GormMembershipRepository) FindByPersonID(personID uint) (*domain.Membership, error) {
	var membership domain.Membership
	err := r.db.Where("person_id = ?", personID).First(&membership).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &membership, nil
}

func (r *GormMembershipRepository) FindByID(id uint) (*domain.Membership, error) {
	var membership domain.Membership
	err := r.db.First(&membership, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &membership, nil
}

func (r *GormMembershipRepository) FindAll() ([]domain.Membership, error) {
	var memberships []domain.Membership
	err := r.db.Find(&memberships).Error
	return memberships, err
}

func (r *GormMembershipRepository) Update(membership *domain.Membership) error {
	return r.db.Save(membership).Error
}

func (r *GormMembershipRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Membership{}, id).Error
}
