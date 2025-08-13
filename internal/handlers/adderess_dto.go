package handlers

// AddressDTO es el DTO para la información de la dirección.
type AddressDTO struct {
	ID      uint   `json:"id,omitempty"`
	Address string `json:"address" binding:"required"`
}
