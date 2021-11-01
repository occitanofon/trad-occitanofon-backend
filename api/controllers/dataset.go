package controllers

import (
	"btradoc/pkg"
	"btradoc/pkg/dataset"
	"btradoc/pkg/dialect"
	"btradoc/pkg/translation"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetDatasets(datasetService dataset.Service, dialectService dialect.Service, translationService translation.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		translatorID := c.Locals("translatorID").(string)

		fullDialectParam := c.Params("full_dialect")
		if len(fullDialectParam) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": ErrDialectNotProvided,
			})
		}

		fullDialect, err := url.QueryUnescape(fullDialectParam)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		if !strings.Contains(fullDialect, "_") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": ErrBadFullDialectFormat,
			})
		}

		occitan := strings.Split(fullDialect, "_")
		dialect := occitan[0]
		subdialect := occitan[1]

		match, err := dialectService.Exists(dialect, subdialect)
		if err != nil {
			logger.Error(err)
			switch e := err.(type) {
			case *pkg.DBError:
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		} else if !match {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": ErrBadFullDialectFormat,
			})
		}

		totalOnGoingTranslations, err := translationService.FetchTotalOnGoingTranslations(fullDialect, translatorID)
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

		// prevent translator to reload then fetch again and again without limit
		if totalOnGoingTranslations > 300 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": ErrTooMuchTranslationsFetched,
			})
		}

		datasets, err := datasetService.FetchByDialect(fullDialect)
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

		// assure that there is still dataset to translate
		if len(datasets) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": ErrNoMoreDataset,
			})
		}

		if err = translationService.AddOnGoingTranslations(fullDialect, translatorID, datasets); err != nil {
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

		return c.JSON(datasets)
	}
}
