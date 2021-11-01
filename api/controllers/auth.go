package controllers

import (
	"btradoc/entities"
	"btradoc/pkg"
	"btradoc/pkg/activity"
	"btradoc/pkg/auth"
	"btradoc/pkg/dialect"
	"fmt"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

const DOMAIN string = "occitanofon.org"

func Login(secretKey string, activityService activity.Service, authService auth.Service, dialectService dialect.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		loginBody := new(entities.LoginBody)
		if err := c.BodyParser(&loginBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		transl, err := authService.Login(strings.TrimSpace(loginBody.Username))
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

		if !transl.Confirmed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": ErrAccountNotConfirmed,
			})
		}

		if transl.Suspended {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": ErrAccountSuspended,
			})
		}

		match, err := argon2id.ComparePasswordAndHash(loginBody.Password, transl.Hpwd)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		} else if !match {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": ErrBadCredentials,
			})
		}

		// permissions are needed in order to access some part of frontend
		if len(transl.Permissions) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": ErrNoPermDialect,
			})
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["aud"] = "https://trad.occitanofon.org"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		claims["iat"] = time.Now().Unix()
		claims["iss"] = "https://api.occitanofon.org"
		claims["sub"] = transl.ID
		claims["dperms"] = transl.CompressPerms()

		accessToken, err := token.SignedString([]byte(secretKey))
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		refreshToken, err := authService.CreateRefreshToken(transl.ID)
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

		activityService.AddOrKeepActive(transl.ID)

		c.Set("Set-Cookie", fmt.Sprintf("refreshToken=%s;Secure;HttpOnly;Path=/;Domain=%s;SameSite=Lax", refreshToken, DOMAIN))

		return c.JSON(accessToken)

	}
}

func RefreshToken(secretKey string, authService auth.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		refreshToken := c.Cookies("refreshToken")

		transl, err := authService.FindByRefreshToken(refreshToken)
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

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["aud"] = "https://trad.occitanofon.org"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		claims["iat"] = time.Now().Unix()
		claims["iss"] = "https://api.occitanofon.org"
		claims["sub"] = transl.ID
		claims["dperms"] = transl.CompressPerms()

		accessToken, err := token.SignedString([]byte(secretKey))
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		c.Set("Set-Cookie", fmt.Sprintf("refreshToken=%s;Secure;HttpOnly;Path=/;Domain=%s;SameSite=Lax", refreshToken, DOMAIN))

		return c.JSON(accessToken)
	}
}

func Logout() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.ClearCookie("refreshToken")

		return c.SendStatus(fiber.StatusOK)
	}
}
