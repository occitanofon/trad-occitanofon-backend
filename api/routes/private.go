package routes

import (
	"btradoc/api/controllers"
	"btradoc/api/middlewares"
	"btradoc/helpers"
	"btradoc/pkg/account"
	"btradoc/pkg/activity"
	"btradoc/pkg/auth"
	"btradoc/pkg/dataset"
	"btradoc/pkg/dialect"
	"btradoc/pkg/mailer"
	"btradoc/pkg/translation"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func PrivateEndpoints(app fiber.Router, services ...interface{}) {
	SECRET_KEY := helpers.GetSigningKey()

	_ = services[0].(account.Service)
	activityService := services[1].(activity.Service)
	_ = services[2].(auth.Service)
	datasetService := services[3].(dataset.Service)
	dialectService := services[4].(dialect.Service)
	_ = services[5].(mailer.Service)
	translationService := services[6].(translation.Service)

	api := app.Group("/p", jwtware.New(jwtware.Config{
		SigningMethod:  jwt.SigningMethodHS256.Name,
		SigningKey:     []byte(SECRET_KEY),
		SuccessHandler: middlewares.JWTSuccessHandler(),
		ErrorHandler:   middlewares.JWTErrorHandler(SECRET_KEY),
	}))
	api.Add("GET", "/occitan", controllers.GetOccitan(dialectService))
	api.Add("GET", "/datasets/:full_dialect", controllers.GetDatasets(datasetService, dialectService, translationService))
	api.Add("POST", "/new_translations", controllers.NewTranslations(datasetService, dialectService, translationService, activityService))
	api.Add("GET", "/online_tanslators", controllers.TotalOnlineTranslators(activityService))
	api.Add("GET", "/logout", controllers.Logout())
}
