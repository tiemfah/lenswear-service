package adminroutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiemfah/lenswear-service/internal/handlers/userhdl"
)

func UserEndPoint(router fiber.Router, hdl *userhdl.HTTPHandler) {
	endpoints := router.Group("")
	{
		endpoints.Post("", hdl.CreateUser)
		endpoints.Get("", hdl.GetUsersAsAdmin)
		endpoints.Put("", hdl.ModifyUser)
		endpoints.Delete("", hdl.ResetUserPassword)
	}
}
