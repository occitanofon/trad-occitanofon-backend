package routes

import (
	"btradoc/api/controllers"
	"btradoc/helpers"
	"btradoc/pkg/account"
	"btradoc/pkg/activity"
	"btradoc/pkg/auth"
	"btradoc/pkg/dataset"
	"btradoc/pkg/dialect"
	"btradoc/pkg/mailer"
	"btradoc/pkg/translation"

	"github.com/gofiber/fiber/v2"
)

func PublicEndpoints(app fiber.Router, services ...interface{}) {
	SECRET_KEY := helpers.GetSigningKey()

	accountService := services[0].(account.Service)
	activityService := services[1].(activity.Service)
	authService := services[2].(auth.Service)
	_ = services[3].(dataset.Service)
	dialectService := services[4].(dialect.Service)
	mailerService := services[5].(mailer.Service)
	translationService := services[6].(translation.Service)

	app.Post("/login", controllers.Login(SECRET_KEY, activityService, authService, dialectService))
	app.Get("/refreshtoken", controllers.RefreshToken(SECRET_KEY, authService))
	app.Post("/register", controllers.Register(accountService))
	app.Post("/send_pwd_reset", controllers.SendPasswordReset(accountService, mailerService))
	app.Get("/secret_questions", controllers.SecretQuestions(accountService))
	app.Get("/confirm_token/:token", controllers.ConfirmPasswordResetToken(accountService))
	app.Post("/update_pwd", controllers.UpdatePassword(accountService))
	app.Get("/translations_files", controllers.TranslationsFiles(translationService))
}
