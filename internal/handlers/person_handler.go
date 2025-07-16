package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
)

type PersonHandler struct {
	personService ports.PersonService
}

func NewPersonHandler(personService ports.PersonService) *PersonHandler {
	return &PersonHandler{personService: personService}
}

type PersonRequest struct {
	Name       string          `json:"name" validate:"required"`
	MiddleName string          `json:"middleName" validate:"required"`
	LastName   string          `json:"lastName" validate:"required"`
	Sex        domain.Sex      `json:"sex" validate:"required,oneof=F M"`
	Birthday   time.Time       `json:"birthday" validate:"required"`
	DocNumber  *string         `json:"docNumber"`
	TypeDoc    *domain.DocType `json:"typeDoc" validate:"omitempty,oneof=DNI CE passport"`
}

type PersonResponse struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	MiddleName string          `json:"middleName"`
	LastName   string          `json:"lastName"`
	Sex        domain.Sex      `json:"sex"`
	Birthday   time.Time       `json:"birthday"`
	DocNumber  *string         `json:"docNumber,omitempty"`
	TypeDoc    *domain.DocType `json:"typeDoc,omitempty"`
}

// CreateOrUpdatePersonForUser godoc
// @Summary Create or update own person information
// @Description Create or update person information for the authenticated user.
// @Tags Person
// @Accept json
// @Produce json
// @Param person body PersonRequest true "Person information"
// @Success 200 {object} PersonResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /protected/person [put]
func (h *PersonHandler) CreateOrUpdatePersonForUser(c *fiber.Ctx) error {
	var req PersonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
	}

	userID, ok := c.Locals("userID").(float64)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "user ID not found in context"})
	}

	userIDUint := uint(userID)
	person := domain.Person{
		UserID:     &userIDUint,
		Name:       req.Name,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		Sex:        req.Sex,
		Birthday:   req.Birthday,
		DocNumber:  req.DocNumber,
		TypeDoc:    req.TypeDoc,
	}

	updatedPerson, err := h.personService.CreateOrUpdatePersonForUser(&person)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(h.mapPersonToResponse(updatedPerson))
}

// CreatePersonByAdmin godoc
// @Summary Create a new person record (Admin)
// @Description Create a new person record, not necessarily linked to a user.
// @Tags Admin
// @Accept json
// @Produce json
// @Param person body PersonRequest true "Person information"
// @Success 201 {object} PersonResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /protected/person [post]
func (h *PersonHandler) CreatePersonByAdmin(c *fiber.Ctx) error {
	var req PersonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "cannot parse JSON"})
	}

	person := &domain.Person{
		Name:       req.Name,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		Sex:        req.Sex,
		Birthday:   req.Birthday,
		DocNumber:  req.DocNumber,
		TypeDoc:    req.TypeDoc,
	}

	createdPerson, err := h.personService.CreatePerson(person)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(h.mapPersonToResponse(createdPerson))
}

// DeletePerson godoc
// @Summary Delete person information
// @Description Delete person information by ID (admin only)
// @Tags Admin
// @Produce json
// @Param id path int true "Person ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /protected/person/{id} [delete]
func (h *PersonHandler) DeletePerson(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "missing person ID"})
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "invalid person ID format"})
	}

	if err := h.personService.DeletePerson(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// SearchPersons godoc
// @Summary Search persons
// @Description Search for persons by a single search term. The search is performed on the full name and document number.
// @Tags Person
// @Produce json
// @Param q query string false "Search term"
// @Success 200 {array} PersonResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /protected/person/search [get]
func (h *PersonHandler) SearchPersons(c *fiber.Ctx) error {
	searchTerm := c.Query("q")

	persons, err := h.personService.SearchPersons(searchTerm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.JSON(h.mapPersonsToResponse(persons))
}

func (h *PersonHandler) mapPersonToResponse(p *domain.Person) PersonResponse {
	return PersonResponse{
		ID:         p.ID,
		Name:       p.Name,
		MiddleName: p.MiddleName,
		LastName:   p.LastName,
		Sex:        p.Sex,
		Birthday:   p.Birthday,
		DocNumber:  p.DocNumber,
		TypeDoc:    p.TypeDoc,
	}
}

func (h *PersonHandler) mapPersonsToResponse(persons []domain.Person) []PersonResponse {
	resp := make([]PersonResponse, len(persons))
	for i, p := range persons {
		resp[i] = h.mapPersonToResponse(&p)
	}
	return resp
}
