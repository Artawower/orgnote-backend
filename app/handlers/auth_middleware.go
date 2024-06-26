package handlers

import (
	"orgnote/app/models"
	"orgnote/app/tools"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var ConfigDefault = Config{
	Filter: nil,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = ConfigDefault.Filter
	}
	return cfg
}

// TODO: master add config for common arrangements
func NewAuthMiddleware() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		rawUser := c.Locals("user")
		if rawUser == (*models.User)(nil) {
			return c.Status(fiber.StatusUnauthorized).JSON(NewHttpError[any](ErrInvalidToken, nil))
		}
		return c.Next()
	}
}

type Config struct {
	Filter       func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	GetUser      func(token string) (*models.User, error)
}

func NewUserInjectMiddleware(config ...Config) func(*fiber.Ctx) error {
	cfg := configDefault(config...)
	if cfg.GetUser == nil {
		log.Fatal().Msg("auth middleware: init new auth middleware: GetUser function is required")
	}

	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}

		token := tools.ExtractBearerTokenFromCtx(c)

		var user *models.User
		var err error

		if token == "" {
			c.Locals("user", user)
			return c.Next()
		}

		user, err = cfg.GetUser(token)
		if err != nil {
			log.Info().Msgf("auth middleware: GetUser: %s", err)
		}

		c.Locals("user", user)
		return c.Next()
	}

}
