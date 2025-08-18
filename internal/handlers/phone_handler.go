package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
)

type PhoneHandler struct {
	phoneService ports.PhoneService
}

func NewPhoneHandler(phoneService ports.PhoneService) *PhoneHandler {
	return &PhoneHandler{phoneService: phoneService}
}

func (h *PhoneHandler) CreateOrUpdatePhone(c *fiber.Ctx) error {
	var req PhoneDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "cannot parse JSON"})
	}

	phone := &domain.Phone{
		ID:       req.ID,
		PersonID: req.PersonID,
		Phone:    req.Phone,
	}

	println("phone", phone.PersonID)
	savedPhone, err := h.phoneService.CreateOrUpdatePhone(phone)

	if err != nil {
		if errors.Is(err, ports.ErrPersonNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "la persona especificada no existe"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	responseDTO := PhoneDTO{
		ID:       savedPhone.ID,
		PersonID: savedPhone.PersonID,
		Phone:    savedPhone.Phone,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": responseDTO})
}

func (h *PhoneHandler) DeletePhone(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "missing phone ID"})
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "invalid phone ID format"})
	}

	userID, ok := c.Locals("userID").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "user ID not found in context"})
	}
	userIDUint := uint(userID)

	if err := h.phoneService.DeletePhoneForUser(uint(id), userIDUint); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
