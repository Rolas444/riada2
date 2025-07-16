package handlers

import (
	"errors"

	"github.com/riada2/config"
	"github.com/riada2/internal/core/ports"
	"github.com/riada2/internal/recaptcha"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler maneja las solicitudes de autenticación.
type AuthHandler struct {
	userService ports.UserService
	cfg         *config.Config
}

// NewAuthHandler crea una nueva instancia de AuthHandler.
func NewAuthHandler(userService ports.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		cfg:         cfg,
	}
}

// LoginRequest define el cuerpo de la solicitud para el endpoint de login.
type LoginRequest struct {
	Username       string `json:"username" example:"testuser"`
	Password       string `json:"password" example:"password123"`
	RecaptchaToken string `json:"recaptchaToken" example:"03AGdBq27..."`
}

// Login godoc
// @Summary      Iniciar sesión de un usuario
// @Description  Inicia sesión con nombre de usuario y contraseña, y devuelve un token JWT. Requiere verificación con reCAPTCHA v3.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "Credenciales de inicio de sesión y token reCAPTCHA"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} ErrorResponse "No se puede procesar el JSON o falta el token reCAPTCHA"
// @Failure      401 {object} ErrorResponse "Credenciales inválidas o fallo en la verificación de reCAPTCHA"
// @Failure      500 {object} ErrorResponse "Error interno del servidor"
// @Router       /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "No se puede procesar el JSON"})
	}

	if req.RecaptchaToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Falta el token reCAPTCHA"})
	}

	// Verificar el token de reCAPTCHA
	recaptchaValid, err := recaptcha.Verify(req.RecaptchaToken, h.cfg.RecaptchaSecretKey)
	if err != nil {
		if errors.Is(err, recaptcha.ErrRecaptchaNotConfigured) {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "reCAPTCHA no está configurado correctamente"})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Error al verificar reCAPTCHA"})
		}
	}

	if !recaptchaValid {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Falló la verificación de reCAPTCHA"})
	}

	// Autenticar usuario usando el servicio
	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Credenciales inválidas"})
	}

	return c.Status(fiber.StatusOK).JSON(LoginResponse{Token: token})
}
