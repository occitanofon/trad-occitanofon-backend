package middlewares

import (
	"btradoc/entities"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const (
	ErrJWTToken string = "Une erreur est survenue lors de l'authentication"
)

func JWTSuccessHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)

		claims := user.Claims.(jwt.MapClaims)
		translatorID := claims["sub"].(string)

		c.Locals("translatorID", translatorID)

		return c.Next()
	}
}

func JWTErrorHandler(secretKey string) func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		bearerJWT := string(c.Request().Header.Peek("Authorization"))

		bearerJWTSplit := strings.Split(bearerJWT, " ")

		if len(bearerJWTSplit) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": ErrJWTToken,
			})
		}

		_, err = jwt.ParseWithClaims(bearerJWTSplit[1], &entities.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err := err.(*jwt.ValidationError); err.Errors == jwt.ValidationErrorExpired {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "TOKEN_EXPIRED",
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": ErrJWTToken,
		})
	}
}
