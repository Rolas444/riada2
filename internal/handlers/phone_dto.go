package handlers

// PhoneDTO es el DTO para la información del teléfono.
type PhoneDTO struct {
	ID    uint   `json:"id,omitempty"`
	Phone string `json:"phone" binding:"required"`
}
