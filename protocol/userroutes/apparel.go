package userroutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiemfah/lenswear-service/internal/handlers/apparelhdl"
)

func ApparelEndPoint(router fiber.Router, hdl *apparelhdl.HTTPHandler) {
	endpoints := router.Group("")
	{
		endpoints.Get("", hdl.GetApparels)
		endpoints.Get(":apparelID", hdl.GetApparelByApparelID)
	}
}
