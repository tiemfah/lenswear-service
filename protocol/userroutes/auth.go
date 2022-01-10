package userroutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiemfah/lenswear-service/internal/handlers/authhdl"
)

func AuthEndPoint(router fiber.Router, hdl *authhdl.HTTPHandler) {
	endpoints := router.Group("")
	{
		endpoints.Post("/login", hdl.Login)
	}
}
