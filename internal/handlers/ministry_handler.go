package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
)

type MinistryHandler struct {
	ministryService       ports.MinistryService
	ministryMemberService ports.MinistryMemberService
}

func NewMinistryHandler(ministryService ports.MinistryService, ministryMemberService ports.MinistryMemberService) *MinistryHandler {
	return &MinistryHandler{ministryService: ministryService, ministryMemberService: ministryMemberService}
}

func (h *MinistryHandler) GetAll(c *fiber.Ctx) error {
	ministries, err := h.ministryService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ministries)
}

// POST /minister
func (h *MinistryHandler) CreateMinistry(c *fiber.Ctx) error {
	var req CreateMinistryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	ministry := &domain.Ministry{
		Name:        req.Name,
		Description: req.Description,
		Mission:     req.Mission,
		Status:      req.Status,
	}
	created, err := h.ministryService.Create(ministry)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

// PUT /ministry
func (h *MinistryHandler) UpdateMinistry(c *fiber.Ctx) error {
	var req UpdateMinistryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	input := &domain.Ministry{
		Name:        req.Name,
		Description: req.Description,
		Mission:     req.Mission,
		Status:      req.Status,
	}
	updated, err := h.ministryService.Update(req.ID, input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updated)
}

// POST /ministry/member
func (h *MinistryHandler) CreateMinistryMember(c *fiber.Ctx) error {
	var req CreateMinistryMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	member := &domain.MinistryMember{
		MinistryID: req.MinistryID,
		PersonID:   req.PersonID,
		Role:       req.Role,
		Status:     req.Status,
	}
	created, err := h.ministryMemberService.Create(member)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

// PUT /ministry/member
func (h *MinistryHandler) UpdateMinistryMember(c *fiber.Ctx) error {
	var req UpdateMinistryMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	input := &domain.MinistryMember{
		MinistryID: req.MinistryID,
		PersonID:   req.PersonID,
		Role:       req.Role,
		Status:     req.Status,
	}
	updated, err := h.ministryMemberService.Update(0, input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updated)
}
