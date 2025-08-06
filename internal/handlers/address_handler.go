package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
)

type AddressHandler struct {
	addressService ports.AddressService
}

func NewAddressHandler(addressService ports.AddressService) *AddressHandler {
	return &AddressHandler{addressService: addressService}
}

func (h *AddressHandler) CreateOrUpdateAddress(c *fiber.Ctx) error {
	var req AddressDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "cannot parse JSON"})
	}

	// El ID del DTO se usará para actualizaciones. Si es 0, es una creación.
	address := &domain.Address{
		ID:      req.ID,
		Address: req.Address,
	}

	userID, ok := c.Locals("userID").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "user ID not found in context"})
	}
	userIDUint := uint(userID)

	savedAddress, err := h.addressService.CreateOrUpdateAddressForUser(address, userIDUint)
	if err != nil {
		// Se podría manejar errores específicos como "no encontrado" o "no autorizado"
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	// Devolver la dirección guardada como un DTO
	responseDTO := AddressDTO{
		ID:      savedAddress.ID,
		Address: savedAddress.Address,
	}

	return c.Status(fiber.StatusOK).JSON(responseDTO)
}

func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "missing address ID"})
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "invalid address ID format"})
	}

	userID, ok := c.Locals("userID").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "user ID not found in context"})
	}
	userIDUint := uint(userID)

	if err := h.addressService.DeleteAddressForUser(uint(id), userIDUint); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
