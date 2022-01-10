package adminroutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiemfah/lenswear-service/internal/handlers/apparelhdl"
)

func ApparelEndPoint(router fiber.Router, hdl *apparelhdl.HTTPHandler) {
	endpoints := router.Group("")
	{
		endpoints.Post("", hdl.CreateApparel)
		endpoints.Get("", hdl.GetApparels)
		endpoints.Delete("", hdl.DeleteApparelByApparelID)
	}
}
