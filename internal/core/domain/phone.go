package domain

import "time"

// Phone representa un teléfono asociado a una persona.
// Corresponde a la tabla 'phones'.
type Phone struct {
	ID        uint
	PersonID  uint
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
