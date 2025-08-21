package domain

import "time"

// Sex define el tipo para el sexo de una persona.
type Sex string

const (
	Female Sex = "F"
	Male   Sex = "M"
)

// DocType define el tipo para el tipo de documento de una persona.
type DocType string

const (
	DNI      DocType = "DNI"
	CE       DocType = "CE"
	Passport DocType = "passport"
)

// Person representa la entidad de una persona en el sistema.
type Person struct {
	ID         uint
	UserID     *uint
	Name       string
	MiddleName string
	LastName   string
	Sex        Sex
	Birthday   *time.Time
	DocNumber  *string
	TypeDoc    *DocType
	Email      *string
	Photo      *string
	Addresses  []Address
	Phones     []Phone
	Membership *Membership // Relación 1:1 con Membership
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
