package main

import (
	"btradoc/api/middlewares"
	"btradoc/api/routes"
	"btradoc/helpers"
	"btradoc/pkg/account"
	"btradoc/pkg/activity"
	"btradoc/pkg/auth"
	"btradoc/pkg/dataset"
	"btradoc/pkg/dialect"
	"btradoc/pkg/mailer"
	"btradoc/pkg/translation"
	"btradoc/storage/mongodb"
	"flag"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

var (
	PRODUCTION_MOD *bool
	LOG_FILE       string
	ALLOW_ORIGINS  string
)

func init() {
	LOG_FILE = "logrus.log"

	PRODUCTION_MOD = flag.Bool("prod", false, "enable production mode")
	flag.Parse()

	if *PRODUCTION_MOD {
		ALLOW_ORIGINS = "https://trad.occitanofon.org"
	} else {
		ALLOW_ORIGINS = "http://127.0.0.1:3333"
	}

	mongodb.InitMongoDatabase()
}

func main() {
	logger := logrus.New()
	logger.ReportCaller = true

	if *PRODUCTION_MOD {
		file := helpers.CreateLogFile(LOG_FILE)
		defer file.Close()

		logger.SetOutput(file)
		logger.Formatter = &logrus.JSONFormatter{}
	} else {
		logger.SetOutput(os.Stdout)
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	db := mongodb.NewMongoClient()

	accountRepo := account.NewRepo(db)
	accountService := account.NewService(accountRepo)

	activityService := activity.NewService()

	authRepo := auth.NewRepo(db)
	authService := auth.NewService(authRepo)

	datasetRepo := dataset.NewRepo(db)
	datasetService := dataset.NewService(datasetRepo)

	dialectRepo := dialect.NewRepo(db)
	dialectService := dialect.NewService(dialectRepo)

	mailerService := mailer.NewService(db, logger)

	translationRepo := translation.NewRepo(db)
	translationService := translation.NewService(translationRepo)

	app := fiber.New()

	if !*PRODUCTION_MOD {
		app.Use(log.New())
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     ALLOW_ORIGINS,
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Accept,Authorization,Content-Type,Cookie,Origin,Set-Cookie",
		AllowCredentials: true,
	}))

	app.Use(middlewares.Logrus(logger))

	routes.PrivateEndpoints(app, accountService, activityService, authService, datasetService, dialectService, mailerService, translationService)
	routes.PublicEndpoints(app, accountService, activityService, authService, datasetService, dialectService, mailerService, translationService)

	_ = app.Listen(":9321")
}
