package domain

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Role es un tipo para los roles de usuario.
type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

// Scan implementa la interfaz Scanner para el tipo Role.
func (r *Role) Scan(value interface{}) error {
	// El driver de la base de datos puede devolver un string o []byte.
	// Manejamos ambos casos para mayor robustez.
	switch v := value.(type) {
	case []byte:
		*r = Role(v)
	case string:
		*r = Role(v)
	default:
		return errors.New("failed to scan role: unsupported type")
	}
	return nil
}

// Value implementa la interfaz driver.Valuer para el tipo Role.
func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

// User representa un usuario en el sistema.
type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         Role   `gorm:"type:varchar(10);not null;default:user"`
	PersonID     *uint
	Person       Person // GORM usará PersonID como clave foránea por convención.
}

// UserResponse es un DTO para enviar datos de usuario sin la contraseña.
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
