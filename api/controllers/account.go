package controllers

import (
	"btradoc/entities"
	"btradoc/helpers"
	"btradoc/pkg"
	"btradoc/pkg/account"
	"btradoc/pkg/mailer"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Register(accountService account.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		registerBody := new(entities.RegisterBody)
		if err := c.BodyParser(&registerBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		registerBody.TrimFields()

		if err := helpers.IsEmailValid(registerBody.Email); err != nil {
			logger.Error(err)
			if e, ok := err.(*helpers.HelperError); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		if err := helpers.UsernameValidity(registerBody.Username); err != nil {
			logger.Error(err)
			if e, ok := err.(*helpers.HelperError); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		if len(registerBody.Password) < 9 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": ErrPasswordTooShort,
			})
		}

		if len(registerBody.SecretQuestionsAndResponses) != 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": ErrSecretQuestions,
			})
		}

		hashedPassword, err := argon2id.CreateHash(registerBody.Password, argon2id.DefaultParams)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		registerBody.SecretQuestionsAndResponses.Trimer()
		registerBody.SecretQuestionsAndResponses.ToLowerResponses()

		// responses are hashed for privacy
		var secretQuestionsAndHashedResponses []entities.SecretQuestionAndResponse

		for _, sqar := range registerBody.SecretQuestionsAndResponses {
			newSQHR := entities.SecretQuestionAndResponse{
				Question: sqar.Question,
			}
			newSQHR.Response, err = argon2id.CreateHash(sqar.Response, argon2id.DefaultParams)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
			secretQuestionsAndHashedResponses = append(secretQuestionsAndHashedResponses, newSQHR)
		}

		newTranslator := entities.NewTranslator{
			Email:                       registerBody.Email,
			Username:                    registerBody.Username,
			Hpwd:                        hashedPassword,
			SecretQuestionsAndResponses: secretQuestionsAndHashedResponses,
		}

		if err = accountService.Create(newTranslator); err != nil {
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

		return c.JSON(fiber.Map{
			"msg": "Tu as bien été inscrit, un administrateur doit valider le compte afin que tu puisses te connecter",
		})
	}
}

func SendPasswordReset(accountService account.Service, mailerService mailer.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		emailBody := new(entities.PasswordResetBody)
		if err := c.BodyParser(&emailBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		if err := helpers.IsEmailValid(emailBody.Email); err != nil {
			logger.Error(err)
			if e, ok := err.(*helpers.HelperError); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		transl, err := accountService.ResetPassword(emailBody.Email)
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

		if allowed := mailerService.IsAllowed(transl.Email); !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": ErrMailerNotAllowed,
			})
		}
		mailerService.SendResetPasswordLink(transl)

		return c.JSON(fiber.Map{
			"msg": "Un email vous sera prochainement envoyé afin de procéder à la réinitialisation de votre mot de passe",
		})
	}
}

func ConfirmPasswordResetToken(accountService account.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		token := c.Params("token")

		translator, err := accountService.FetchSecretQuestionsAndResponses(token)
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

		// Pull out responses only
		var secretQuestionsOnly []string
		for question := range translator.SecretQuestionsAndResponses {
			secretQuestionsOnly = append(secretQuestionsOnly, question)
		}

		return c.JSON(secretQuestionsOnly)
	}
}

func SecretQuestions(accountService account.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		secretQuestions, err := accountService.FetchAllSecretQuestions()
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

		return c.JSON(secretQuestions)
	}
}

func UpdatePassword(accountService account.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		updatePasswordBody := new(entities.PasswordUpdateBody)
		if err := c.BodyParser(&updatePasswordBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		updatePasswordBody.SecretQuestionsAndResponses.Trimer()
		updatePasswordBody.SecretQuestionsAndResponses.ToLowerResponses()

		result, err := accountService.FetchSecretQuestionsAndResponses(updatePasswordBody.Token)
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

		for _, questionAndResponse := range updatePasswordBody.SecretQuestionsAndResponses {
			question := questionAndResponse.Question
			response := questionAndResponse.Response

			translHashedResponse, has := result.SecretQuestionsAndResponses[question]
			if !has {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": ErrSecretQuestionsNoMatch,
				})
			}

			match, _, err := argon2id.CheckHash(response, translHashedResponse)
			if !match {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": ErrSecretQuestionsNoMatch,
				})
			} else if err != nil {
				logger.Error(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		}

		newHashedPassword, err := argon2id.CreateHash(strings.TrimSpace(updatePasswordBody.Password), argon2id.DefaultParams)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		if err = accountService.UpdatePassword(result.TranslatorID.Hex(), newHashedPassword); err != nil {
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

		// remove temporary token

		return c.JSON(fiber.Map{
			"msg": "Votre mot de passe a bien été modifié",
		})
	}
}
