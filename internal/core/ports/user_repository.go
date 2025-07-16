package ports

import "github.com/riada2/internal/core/domain"

// UserRepository es el puerto para la persistencia de usuarios.
type UserRepository interface {
	Save(user *domain.User) error
	FindByUsername(username string) (*domain.User, error)
}
