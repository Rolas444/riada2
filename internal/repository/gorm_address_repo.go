package repository

import (
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type gormAddressRepository struct {
	db *gorm.DB
}

func NewGormAddressRepository(db *gorm.DB) ports.AddressRepository {
	return &gormAddressRepository{db: db}
}

func (r *gormAddressRepository) Save(address *domain.Address) error {
	return r.db.Save(address).Error
}

func (r *gormAddressRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Address{}, id).Error
}

func (r *gormAddressRepository) FindByID(id uint) (*domain.Address, error) {
	var address domain.Address
	if err := r.db.First(&address, id).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *gormAddressRepository) CountByPersonID(personID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Address{}).Where("person_id = ?", personID).Count(&count).Error
	return count, err
}
