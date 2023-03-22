package handlers

import (
	"moonbrain/app/configs"
	_ "moonbrain/app/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
		"moonbrain/app/docs"

)

func RegisterSwagger(app fiber.Router, config configs.Config) {
	docs.SwaggerInfo.Host = config.BackendDomain
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.BasePath = "/v1"


	app.Get("/swagger/*", swagger.HandlerDefault) // default
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "doc.json",
		DeepLinking: false,
		DocExpansion: "none",
	}))
}
