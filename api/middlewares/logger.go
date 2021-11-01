package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Logrus(logger *logrus.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals("logger", logger)
		return c.Next()
	}
}
