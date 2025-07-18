package handlers

import domain "github.com/riada2/internal/core/domain"

// LoginResponseUser contiene los detalles del usuario para la respuesta de login.
type LoginResponseUser struct {
	Email string       `json:"email" example:"testuser"`
	Role  *domain.Role `json:"role" example:"user"`
}

// LoginResponse representa el cuerpo de la respuesta para el endpoint de login.
type LoginResponse struct {
	Token string            `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  LoginResponseUser `json:"user"`
}
