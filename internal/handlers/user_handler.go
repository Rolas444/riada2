package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/riada2/internal/core/ports"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type RegisterRequest struct {
	Username string `json:"username" example:"newuser"`
	Password string `json:"password" example:"password123"`
}

// RegisterResponse define la estructura de la respuesta para un registro exitoso.
type RegisterResponse struct {
	ID        uint      `json:"id" example:"1"`
	Username  string    `json:"username" example:"newuser"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"createdAt"`
}

// LoginResponse define la estructura de la respuesta para un login exitoso.
// type LoginResponse struct {
// 	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
// }

// ProfileResponse define la estructura de la respuesta para el perfil de usuario.
type ProfileResponse struct {
	Message string `json:"message" example:"Welcome!"`
	UserID  uint   `json:"userID" example:"1"`
	Role    string `json:"role" example:"user"`
}

// ErrorResponse define una estructura estándar para los errores.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Register godoc
// @Summary      Registrar un nuevo usuario
// @Description  Crea una nueva cuenta de usuario con el nombre de usuario y contraseña proporcionados.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body RegisterRequest true "Información de registro del usuario"
// @Success      201 {object} RegisterResponse
// @Failure      400 {object} ErrorResponse "No se puede procesar el JSON"
// @Failure      409 {object} ErrorResponse "El nombre de usuario ya existe"
// @Router       /register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "cannot parse JSON"})
	}

	user, err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(ErrorResponse{Error: err.Error()})
	}

	response := RegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetProfile godoc
// @Summary      Obtener perfil de usuario
// @Description  Obtiene la información del perfil del usuario autenticado actualmente. Requiere token JWT.
// @Tags         User
// @Produce      json
// @Success      200 {object} ProfileResponse
// @Failure      401 {object} ErrorResponse "No autorizado"
// @Security     ApiKeyAuth
// @Router       /protected/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// La información del usuario se inyecta desde el middleware en c.Locals
	// El ID de usuario del token JWT (claim "sub") se decodifica como float64 por defecto en Go.
	userIDClaim, _ := c.Locals("userID").(float64)
	userRoleClaim, _ := c.Locals("userRole").(string)

	response := ProfileResponse{
		Message: "Welcome!",
		UserID:  uint(userIDClaim),
		Role:    userRoleClaim,
	}
	return c.JSON(response)
}
