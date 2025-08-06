package repository

import (
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type gormPhoneRepository struct {
	db *gorm.DB
}

func NewGormPhoneRepository(db *gorm.DB) ports.PhoneRepository {
	return &gormPhoneRepository{db: db}
}

func (r *gormPhoneRepository) Save(phone *domain.Phone) error {
	return r.db.Save(phone).Error
}

func (r *gormPhoneRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Phone{}, id).Error
}

func (r *gormPhoneRepository) FindByID(id uint) (*domain.Phone, error) {
	var phone domain.Phone
	if err := r.db.First(&phone, id).Error; err != nil {
		return nil, err
	}
	return &phone, nil
}

func (r *gormPhoneRepository) CountByPersonID(personID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Phone{}).Where("person_id = ?", personID).Count(&count).Error
	return count, err
}
