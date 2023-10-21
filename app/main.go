package main

import (
	"context"
	"net/http"
	"orgnote/app/configs"
	"orgnote/app/handlers"
	"orgnote/app/infrastructure"
	"orgnote/app/repositories"
	"orgnote/app/services"
	"os"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Org Note API
// @version 0.0.1
// @description List of methods for work with Org Note.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email artawower@protonmail.com
// @license.name GPL 3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
func main() {
	// TODO: master use DIG for dependencies
	config := configs.NewConfig()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	http := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongo")
		return
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to ping mongo: %v", err)
		return
	}

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	database := mongoClient.Database("orgnote")

	subscriptionAPI, err := infrastructure.NewSubscription(
		http,
		config.AccessCheckerURL,
		config.AccessCheckToken,
		cache.New[string, infrastructure.SubscriptionInfo],
		config.AccessTokenCacheLifeTime,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create subscription")
		return
	}

	app := fiber.New(fiber.Config{
		BodyLimit: config.MaximumFileSize,
	})
	api := app.Group("/v1")

	// TODO: master May be someday there will be DI
	noteRepository := repositories.NewNoteRepository(database)
	tagRepository := repositories.NewTagRepository(database)
	userRepository := repositories.NewUserRepository(database)
	fileStorage := infrastructure.NewFileStorage(config.MediaPath)

	app.Use(cors.New())
	app.Use(handlers.NewUserInjectMiddleware(handlers.Config{
		GetUser: userRepository.FindUserByToken,
	}))

	authMiddleware := handlers.NewAuthMiddleware()
	accessMiddleware := handlers.NewAccessMiddleware(subscriptionAPI)

	noteService := services.NewNoteService(noteRepository, userRepository, tagRepository, fileStorage)
	tagService := services.NewTagService(tagRepository)
	userService := services.NewUserService(userRepository, noteRepository, subscriptionAPI)
	fileService := services.NewFileService(fileStorage, userRepository)

	// api.Use(handlers.NewAuthMiddleware())
	// TODO: expose to external fn

	handlers.RegisterSwagger(api, config)
	handlers.RegisterNoteHandler(api, noteService, authMiddleware, accessMiddleware)
	handlers.RegisterTagHandler(api, tagService)
	handlers.RegisterAuthHandler(api, userService, config, authMiddleware)
	handlers.RegisterFileHandler(api, fileService, authMiddleware, accessMiddleware)
	// handlers.RegisterUserHandlers(app)
	// handlers.RegisterTagHandlers(app)
	app.Static("media", config.MediaPath)

	// NOTE: for local file uploading (tmp quick hack)
	if config.Debug {
		app.Static("v1/media", config.MediaPath)
	}

	log.Info().Msg("Application start debug mode: " + config.AppAddress)
	app.Listen(config.AppAddress)
}
