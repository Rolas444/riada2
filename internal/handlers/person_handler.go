package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/riada2/internal/core/ports"
)

type PersonHandler struct {
	personService ports.PersonService
}

func NewPersonHandler(personService ports.PersonService) *PersonHandler {
	return &PersonHandler{personService: personService}
}

// CreateOrUpdatePersonForUser godoc
// @Summary Create or update own person information
// @Description Create or update person information for the authenticated user.
// @Tags Person
// @Accept json
// @Produce json
// @Param person body handlers.PersonRequest true "Person information"
// @Success 200 {object} handlers.PersonResponse
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

	person, err := req.ToDomain()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "invalid birthday format, use YYYY-MM-DD"})
	}

	userID, ok := c.Locals("userID").(float64)
	if !ok {
		// This should ideally not happen if the auth middleware is working correctly.
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "user ID not found in context"})
	}

	userIDUint := uint(userID)
	person.UserID = &userIDUint

	updatedPerson, err := h.personService.CreateOrUpdatePersonForUser(person)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(NewPersonResponse(updatedPerson))
}

// CreatePersonByAdmin godoc
// @Summary Create a new person record (Admin)
// @Description Create a new person record, not necessarily linked to a user.
// @Tags Admin
// @Accept json
// @Produce json
// @Param person body handlers.PersonRequest true "Person information"
// @Success 201 {object} handlers.PersonResponse
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

	person, err := req.ToDomain()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "invalid birthday format, use YYYY-MM-DD"})
	}

	createdPerson, err := h.personService.CreatePerson(person)
	if err != nil {
		// Consider more specific error codes, e.g., 409 Conflict if person exists.
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(NewPersonResponse(createdPerson))
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
		if errors.Is(err, ports.ErrPersonNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: err.Error()})
		}
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
// @Success 200 {array} handlers.PersonResponse
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

	// Convert domain objects to response DTOs
	responseDTOs := make([]PersonResponse, len(persons))
	for i, p := range persons {
		responseDTOs[i] = NewPersonResponse(&p)
	}

	return c.JSON(responseDTOs)
}
