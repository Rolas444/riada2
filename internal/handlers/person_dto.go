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
	Addresses  []AddressDTO    `json:"addresses,omitempty"`
	Phones     []PhoneDTO      `json:"phones,omitempty"`
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

	// Convertir DTOs de dirección a modelos de dominio
	var addresses []domain.Address
	if pr.Addresses != nil {
		for _, addrDTO := range pr.Addresses {
			addresses = append(addresses, domain.Address{
				ID:      addrDTO.ID,
				Address: addrDTO.Address,
			})
		}
	}

	// Convertir DTOs de teléfono a modelos de dominio
	var phones []domain.Phone
	if pr.Phones != nil {
		for _, phoneDTO := range pr.Phones {
			phones = append(phones, domain.Phone{
				ID:    phoneDTO.ID,
				Phone: phoneDTO.Phone,
			})
		}
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
		Addresses:  addresses,
		Phones:     phones,
	}, nil
}

// MembershipDTO es el DTO para la información de membresía en las respuestas de Person
type MembershipDTO struct {
	ID               uint       `json:"id,omitempty"`
	StartedAt        *string    `json:"startedAt,omitempty"` // Se envía en formato "YYYY-MM-DD"
	MembershipSigned bool       `json:"membershipSigned"`
	State            string     `json:"state"`
	Transferred      bool       `json:"transferred"`
	NameLastChurch   *string    `json:"nameLastChurch,omitempty"`
	Baptized         bool       `json:"baptized"`
	BaptismDate      *string    `json:"baptismDate,omitempty"` // Se envía en formato "YYYY-MM-DD"
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
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
	Addresses  []AddressDTO    `json:"addresses,omitempty"`
	Phones     []PhoneDTO      `json:"phones,omitempty"`
	Membership *MembershipDTO  `json:"membership,omitempty"`
}

// NewPersonResponse es una función constructora que convierte una entidad
// domain.Person a un DTO PersonResponse, asegurando que el formato de los datos
// sea el correcto para la API.
func NewPersonResponse(person *domain.Person) PersonResponse {
	var birthdayStr string
	if person.Birthday != nil {
		birthdayStr = person.Birthday.Format("2006-01-02")
	}

	// Convertir modelos de dominio de dirección a DTOs
	var addressDTOs []AddressDTO
	if person.Addresses != nil {
		for _, addr := range person.Addresses {
			addressDTOs = append(addressDTOs, AddressDTO{
				ID:       addr.ID,
				PersonID: person.ID,
				Address:  addr.Address,
			})
		}
	}

	// Convertir modelos de dominio de teléfono a DTOs
	var phoneDTOs []PhoneDTO
	if person.Phones != nil {
		for _, phone := range person.Phones {
			phoneDTOs = append(phoneDTOs, PhoneDTO{
				ID:       phone.ID,
				PersonID: person.ID,
				Phone:    phone.Phone,
			})
		}
	}

	// Convertir modelo de dominio de membresía a DTO
	var membershipDTO *MembershipDTO
	if person.Membership != nil {
		var startedAtStr *string
		if person.Membership.StartedAt != nil {
			startedAt := person.Membership.StartedAt.Format("2006-01-02")
			startedAtStr = &startedAt
		}

		var baptismDateStr *string
		if person.Membership.BaptismDate != nil {
			baptismDate := person.Membership.BaptismDate.Format("2006-01-02")
			baptismDateStr = &baptismDate
		}

		membershipDTO = &MembershipDTO{
			ID:               person.Membership.ID,
			StartedAt:        startedAtStr,
			MembershipSigned: person.Membership.MembershipSigned,
			State:            string(person.Membership.State),
			Transferred:      person.Membership.Transferred,
			NameLastChurch:   person.Membership.NameLastChurch,
			Baptized:         person.Membership.Baptized,
			BaptismDate:      baptismDateStr,
			CreatedAt:        &person.Membership.CreatedAt,
			UpdatedAt:        &person.Membership.UpdatedAt,
		}
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
		Addresses:  addressDTOs,
		Phones:     phoneDTOs,
		Membership: membershipDTO,
	}
}
