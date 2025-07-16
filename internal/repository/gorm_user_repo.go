package repository

import (
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) ports.UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Save(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *gormUserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
