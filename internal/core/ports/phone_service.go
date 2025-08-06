package ports

import "github.com/riada2/internal/core/domain"

type PhoneService interface {
	CreateOrUpdatePhoneForUser(phone *domain.Phone, userID uint) (*domain.Phone, error)
	DeletePhoneForUser(phoneID uint, userID uint) error
}
