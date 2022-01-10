package middlewares

import (
	"crypto/rsa"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func AuthMiddleware(publicKey *rsa.PublicKey) fiber.Handler {
	return jwtware.New(jwtware.Config{
		AuthScheme:    "Bearer",
		SigningMethod: "RS256",
		SigningKey:    publicKey,
		ErrorHandler:  jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "missing or invalid token"})
}
