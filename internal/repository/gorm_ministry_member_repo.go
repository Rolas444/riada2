package repository

import (
	"errors"

	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type GormMinistryMemberRepository struct {
	db *gorm.DB
}

func NewGormMinistryMemberRepository(db *gorm.DB) ports.MinistryMemberRepository {
	return &GormMinistryMemberRepository{db: db}
}

func (r *GormMinistryMemberRepository) Save(member *domain.MinistryMember) error {
	return r.db.Create(member).Error
}

func (r *GormMinistryMemberRepository) Update(member *domain.MinistryMember) error {
	return r.db.Save(member).Error
}

func (r *GormMinistryMemberRepository) FindByMinistryAndPerson(ministryID, personID uint) (*domain.MinistryMember, error) {
	var mm domain.MinistryMember
	err := r.db.Where("ministry_id = ? AND person_id = ?", ministryID, personID).First(&mm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mm, nil
}
