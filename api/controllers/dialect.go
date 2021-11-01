package controllers

import (
	"btradoc/pkg"
	"btradoc/pkg/dialect"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetOccitan(dialectService dialect.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		translatorID := c.Locals("translatorID").(string)

		occitan, err := dialectService.FetchOccitan(translatorID)
		if err != nil {
			logger.Error(err)
			if e, ok := err.(*pkg.DBError); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		return c.JSON(occitan)
	}
}
