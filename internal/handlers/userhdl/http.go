package userhdl

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/internal/core/ports"
)

type HTTPHandler struct {
	userService ports.UserService
}

func NewHTTPHandler(userService ports.UserService) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
	}
}

func (h *HTTPHandler) CreateUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	in := new(domain.CreateUserReq)
	in.Requester = &domain.Requester{
		UserID:   user["UserID"].(string),
		UserRole: user["UserRole"].(string),
	}
	if err := c.BodyParser(in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := in.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.userService.CreateUser(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *HTTPHandler) GetUsersAsAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	in := new(domain.GetUsersAsAdminReq)
	in.Limit = limit
	in.Offset = offset
	in.Requester = &domain.Requester{
		UserID:   user["UserID"].(string),
		UserRole: user["UserRole"].(string),
	}
	if err := in.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.userService.GetUsersAsAdmin(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *HTTPHandler) Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	res, err := h.userService.GetUserByUserID(c.Context(), &domain.GetUserByUserIDReq{UserID: user["UserID"].(string)})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *HTTPHandler) ModifyUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	in := new(domain.ModifyUserReq)
	in.Requester = &domain.Requester{
		UserID:   user["UserID"].(string),
		UserRole: user["UserRole"].(string),
	}
	if err := c.BodyParser(in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := in.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.userService.ModifyUser(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *HTTPHandler) ResetUserPassword(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	in := new(domain.ResetUserPasswordReq)
	in.Requester = &domain.Requester{
		UserID:   user["UserID"].(string),
		UserRole: user["UserRole"].(string),
	}
	if err := c.BodyParser(in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := in.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.userService.ResetUserPassword(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
