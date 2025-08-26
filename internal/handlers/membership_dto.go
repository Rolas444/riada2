package handlers

import "time"

// CreateMembershipRequest define la estructura de la solicitud para crear una membresía.
type CreateMembershipRequest struct {
	PersonID         uint       `json:"personID" example:"1" validate:"required"`
	StartedAt        *time.Time `json:"startedAt" example:"2024-01-01T00:00:00Z"`
	MembershipSigned bool       `json:"membershipSigned" example:"false"`
	State            string     `json:"state" example:"A" validate:"omitempty,oneof=A I O S"`
	Transferred      bool       `json:"transferred" example:"false"`
	NameLastChurch   *string    `json:"nameLastChurch" example:"Iglesia Anterior"`
	Baptized         bool       `json:"baptized" example:"false"`
	BaptismDate      *time.Time `json:"baptismDate" example:"2024-01-01T00:00:00Z"`
}

// MembershipResponse define la estructura de la respuesta para una membresía.
type MembershipResponse struct {
	ID               uint       `json:"id" example:"1"`
	PersonID         uint       `json:"personID" example:"1"`
	StartedAt        *time.Time `json:"startedAt" example:"2024-01-01T00:00:00Z"`
	MembershipSigned bool       `json:"membershipSigned" example:"false"`
	State            string     `json:"state" example:"A"`
	Transferred      bool       `json:"transferred" example:"false"`
	NameLastChurch   *string    `json:"nameLastChurch" example:"Iglesia Anterior"`
	Baptized         bool       `json:"baptized" example:"false"`
	BaptismDate      *time.Time `json:"baptismDate" example:"2024-01-01T00:00:00Z"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

// UpdateMembershipRequest define la estructura de la solicitud para actualizar una membresía.
type UpdateMembershipRequest struct {
	ID               uint       `json:"id" example:"1" validate:"required"`
	PersonID         uint       `json:"personID" example:"1" validate:"required"`
	StartedAt        *time.Time `json:"startedAt" example:"2024-01-01T00:00:00Z"`
	MembershipSigned bool       `json:"membershipSigned" example:"false"`
	State            string     `json:"state" example:"A" validate:"omitempty,oneof=A I O S"`
	Transferred      bool       `json:"transferred" example:"false"`
	NameLastChurch   *string    `json:"nameLastChurch" example:"Iglesia Anterior"`
	Baptized         bool       `json:"baptized" example:"false"`
	BaptismDate      *time.Time `json:"baptismDate" example:"2024-01-01T00:00:00Z"`
}
