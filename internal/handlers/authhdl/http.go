package authhdl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/internal/core/ports"
)

type HTTPHandler struct {
	authService ports.AuthenticationService
}

func NewHTTPHandler(authService ports.AuthenticationService) *HTTPHandler {
	return &HTTPHandler{
		authService: authService,
	}
}

func (h *HTTPHandler) Login(c *fiber.Ctx) error {
	in := new(domain.LoginReq)
	if err := c.BodyParser(in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := in.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.authService.Login(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func Test(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
