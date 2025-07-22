package ports

import "github.com/riada2/internal/core/domain"

// UserService es el puerto para la lógica de negocio de usuarios.
type UserService interface {
	Register(username, password string) (*domain.User, error)
	Login(username, password string) (string, *domain.Role, error) // Devuelve el token JWT

	// GetAllUsers devuelve una lista de todos los usuarios sin sus contraseñas.
	GetAllUsers() ([]domain.UserResponse, error)
	// UpdateUser actualiza la información de un usuario (e.g., username, role).
	UpdateUser(id uint, username *string, role *domain.Role) (*domain.UserResponse, error)
}
