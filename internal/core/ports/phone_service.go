package ports

import "github.com/riada2/internal/core/domain"

type PhoneService interface {
	CreateOrUpdatePhone(phone *domain.Phone) (*domain.Phone, error)
	DeletePhoneForUser(phoneID uint, userID uint) error
}
