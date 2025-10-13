package repository

import (
	"errors"

	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type GormMinistryRepository struct {
	db *gorm.DB
}

func NewGormMinistryRepository(db *gorm.DB) ports.MinistryRepository {
	return &GormMinistryRepository{db: db}
}

func (r *GormMinistryRepository) Save(ministry *domain.Ministry) error {
	return r.db.Create(ministry).Error
}

func (r *GormMinistryRepository) Update(ministry *domain.Ministry) error {
	return r.db.Save(ministry).Error
}

func (r *GormMinistryRepository) FindByID(id uint) (*domain.Ministry, error) {
	var m domain.Ministry
	err := r.db.First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
