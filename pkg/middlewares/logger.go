package middlewares

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoggerMiddleware() func(*fiber.Ctx) error {
	cfg := logger.Config{
		Next:         nil,
		Format:       "[${time}] | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stderr,
	}
	err := logger.New(cfg)
	return err
}
