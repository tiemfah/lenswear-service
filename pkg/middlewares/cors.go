package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware() func(*fiber.Ctx) error {
	cfg := cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Accept-Language, Authorization, X-Requested-With",
	}
	err := cors.New(cfg)
	return err
}
