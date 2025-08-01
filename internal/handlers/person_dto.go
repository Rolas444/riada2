package handlers

import (
	"time"

	"github.com/riada2/internal/core/domain"
)

// PersonRequest es el DTO (Data Transfer Object) para recibir los datos
// de una persona en las solicitudes HTTP (creación o actualización).
// La validación de los campos (ej. que no estén vacíos) se realiza en el handler.
type PersonRequest struct {
	ID         *uint           `json:"id"`
	Name       string          `json:"name"`
	MiddleName string          `json:"middleName"`
	LastName   string          `json:"lastName"`
	Sex        domain.Sex      `json:"sex"`
	Birthday   *string         `json:"birthday,omitempty"` // Se espera el formato "YYYY-MM-DD"
	DocNumber  *string         `json:"docNumber,omitempty"`
	TypeDoc    *domain.DocType `json:"typeDoc,omitempty"`
	Email      *string         `json:"email,omitempty"`
	Photo      *string         `json:"photo,omitempty"`
}

// ToDomain convierte el DTO PersonRequest a la entidad del dominio domain.Person.
// Realiza la conversión de tipos necesarios, como parsear la fecha de nacimiento.
func (pr *PersonRequest) ToDomain() (*domain.Person, error) {
	var birthday *time.Time
	if pr.Birthday != nil && *pr.Birthday != "" {
		parsedBirthday, err := time.Parse("2006-01-02", *pr.Birthday)
		if err != nil {
			return nil, err
		}
		birthday = &parsedBirthday
	}

	var id uint
	if pr.ID != nil {
		id = *pr.ID
	}

	return &domain.Person{
		ID:         id,
		Name:       pr.Name,
		MiddleName: pr.MiddleName,
		LastName:   pr.LastName,
		Sex:        pr.Sex,
		Birthday:   birthday,
		DocNumber:  pr.DocNumber,
		TypeDoc:    pr.TypeDoc,
		Email:      pr.Email,
		Photo:      pr.Photo,
	}, nil
}

// PersonResponse es el DTO para enviar la información de una persona en las respuestas HTTP.
type PersonResponse struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	MiddleName string          `json:"middleName"`
	LastName   string          `json:"lastName"`
	Sex        domain.Sex      `json:"sex"`
	Birthday   string          `json:"birthday,omitempty"` // Se envía en formato "YYYY-MM-DD"
	DocNumber  *string         `json:"docNumber,omitempty"`
	TypeDoc    *domain.DocType `json:"typeDoc,omitempty"`
	Email      *string         `json:"email,omitempty"`
	Photo      *string         `json:"photo,omitempty"`
}

// NewPersonResponse es una función constructora que convierte una entidad
// domain.Person a un DTO PersonResponse, asegurando que el formato de los datos
// sea el correcto para la API.
func NewPersonResponse(person *domain.Person) PersonResponse {
	var birthdayStr string
	if person.Birthday != nil {
		birthdayStr = person.Birthday.Format("2006-01-02")
	}

	return PersonResponse{
		ID:         person.ID,
		Name:       person.Name,
		MiddleName: person.MiddleName,
		LastName:   person.LastName,
		Sex:        person.Sex,
		Birthday:   birthdayStr,
		DocNumber:  person.DocNumber,
		TypeDoc:    person.TypeDoc,
		Email:      person.Email,
		Photo:      person.Photo,
	}
}
