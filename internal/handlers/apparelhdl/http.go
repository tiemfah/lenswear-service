package apparelhdl

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/internal/core/ports"
)

type HTTPHandler struct {
	apparelService ports.ApparelService
}

func NewHTTPHandler(apparelService ports.ApparelService) *HTTPHandler {
	return &HTTPHandler{
		apparelService: apparelService,
	}
}

func (h *HTTPHandler) CreateApparel(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	in := &domain.CreateApparelReq{
		Requester: &domain.Requester{
			UserID:   user["UserID"].(string),
			UserRole: user["UserRole"].(string),
		},
		ApparelTypeID: c.FormValue("ApparelTypeID", ""),
		Name:          c.FormValue("Name", ""),
		Brand:         c.FormValue("Brand", ""),
		Price:         c.FormValue("Price", ""),
		StoreURL:      c.FormValue("StoreURL", ""),
		Files:         form,
	}
	if err := in.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.apparelService.CreateApparel(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *HTTPHandler) GetApparels(c *fiber.Ctx) error {
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	in := new(domain.GetApparelsReq)
	in.Limit = limit
	in.Offset = offset
	if err := c.BodyParser(in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := in.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.apparelService.GetApparels(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *HTTPHandler) GetApparelByApparelID(c *fiber.Ctx) error {
	in := new(domain.GetApparelByApparelIDReq)
	if err := c.BodyParser(in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := in.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.apparelService.GetApparelByApparelID(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *HTTPHandler) DeleteApparelByApparelID(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	in := &domain.DeleteApparelByApparelIDReq{
		Requester: &domain.Requester{
			UserID:   user["UserID"].(string),
			UserRole: user["UserRole"].(string),
		},
		ApparelID: c.Params("apparelID", ""),
	}
	if err := in.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.apparelService.DeleteApparelByApparelID(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
