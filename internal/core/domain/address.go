package domain

import "time"

// Address representa una dirección asociada a una persona.
// Corresponde a la tabla 'addresses'.
type Address struct {
	ID        uint
	PersonID  uint
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
