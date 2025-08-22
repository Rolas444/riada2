package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
)

type MembershipHandler struct {
	membershipService ports.MembershipService
	personService     ports.PersonService
}

func NewMembershipHandler(membershipService ports.MembershipService, personService ports.PersonService) *MembershipHandler {
	return &MembershipHandler{
		membershipService: membershipService,
		personService:     personService,
	}
}

// CreateMembership godoc
// @Summary      Crear una nueva membresía
// @Description  Crea una nueva membresía para una persona. Solo usuarios autenticados pueden crear membresías.
// @Tags         Membership
// @Accept       json
// @Produce      json
// @Param        membership body CreateMembershipRequest true "Información de la membresía"
// @Success      201 {object} MembershipResponse
// @Failure      400 {object} ErrorResponse "Datos inválidos o persona no encontrada"
// @Failure      401 {object} ErrorResponse "No autorizado"
// @Failure      409 {object} ErrorResponse "La persona ya tiene una membresía"
// @Failure      500 {object} ErrorResponse "Error interno del servidor"
// @Security     ApiKeyAuth
// @Router       /protected/membership [post]
func (h *MembershipHandler) CreateMembership(c *fiber.Ctx) error {
	var req CreateMembershipRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Verificar que la persona existe
	_, err := h.personService.GetPersonByID(req.PersonID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "person not found"})
	}

	// Verificar que la persona no tenga ya una membresía
	existingMembership, err := h.membershipService.GetMembershipByPersonID(req.PersonID)
	if err == nil && existingMembership != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "person already has a membership"})
	}

	// Crear la membresía
	membership := &domain.Membership{
		PersonID:         req.PersonID,
		StartedAt:        req.StartedAt,
		MembershipSigned: req.MembershipSigned,
		State:            domain.MembershipState(req.State),
		Transferred:      req.Transferred,
		NameLastChurch:   req.NameLastChurch,
		Baptized:         req.Baptized,
		BaptismDate:      req.BaptismDate,
	}

	// Si no se especifica el estado, establecer como activo por defecto
	if req.State == "" {
		membership.State = domain.Active
	}

	createdMembership, err := h.membershipService.CreateMembership(req.PersonID, membership)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create membership"})
	}

	response := MembershipResponse{
		ID:               createdMembership.ID,
		PersonID:         createdMembership.PersonID,
		StartedAt:        createdMembership.StartedAt,
		MembershipSigned: createdMembership.MembershipSigned,
		State:            string(createdMembership.State),
		Transferred:      createdMembership.Transferred,
		NameLastChurch:   createdMembership.NameLastChurch,
		Baptized:         createdMembership.Baptized,
		BaptismDate:      createdMembership.BaptismDate,
		CreatedAt:        createdMembership.CreatedAt,
		UpdatedAt:        createdMembership.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetMembershipByPersonID godoc
// @Summary      Obtener membresía por ID de persona
// @Description  Obtiene la membresía de una persona específica.
// @Tags         Membership
// @Produce      json
// @Param        personID path int true "ID de la persona"
// @Success      200 {object} MembershipResponse
// @Failure      401 {object} ErrorResponse "No autorizado"
// @Failure      404 {object} ErrorResponse "Membresía no encontrada"
// @Security     ApiKeyAuth
// @Router       /protected/membership/person/{personID} [get]
func (h *MembershipHandler) GetMembershipByPersonID(c *fiber.Ctx) error {
	personIDStr := c.Params("personID")
	personID, err := strconv.ParseUint(personIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid person ID"})
	}

	membership, err := h.membershipService.GetMembershipByPersonID(uint(personID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "membership not found"})
	}

	response := MembershipResponse{
		ID:               membership.ID,
		PersonID:         membership.PersonID,
		StartedAt:        membership.StartedAt,
		MembershipSigned: membership.MembershipSigned,
		State:            string(membership.State),
		Transferred:      membership.Transferred,
		NameLastChurch:   membership.NameLastChurch,
		Baptized:         membership.Baptized,
		BaptismDate:      membership.BaptismDate,
		CreatedAt:        membership.CreatedAt,
		UpdatedAt:        membership.UpdatedAt,
	}

	return c.JSON(response)
}
