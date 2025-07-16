package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// WelcomeResponse define la estructura de la respuesta para la ruta de bienvenida.
type WelcomeResponse struct {
	Message string `json:"message" example:"Welcome to the Riada2 API!"`
	Docs    string `json:"docs" example:"/swagger/index.html"`
}

// Welcome godoc
// @Summary      Muestra un mensaje de bienvenida
// @Description  Obtiene un mensaje de bienvenida de la API con un enlace a la documentación.
// @Tags         General
// @Produce      json
// @Success      200  {object}  WelcomeResponse
// @Router       / [get]
func Welcome(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(WelcomeResponse{
		Message: "Welcome to the Riada2 API!",
		Docs:    "/swagger/index.html",
	})
}
