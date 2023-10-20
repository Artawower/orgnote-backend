package handlers

import (
	"orgnote/app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Subscription interface {
	Check(userEmail string, occupiedSpace int64, err chan<- error)
}

func NewAccessMiddleware(subscription Subscription) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)

		if user == (*models.User)(nil) {
			return c.Status(fiber.StatusBadRequest).JSON(NewHttpError[any](ErrAuthRequired, nil))
		}

		if user.Email == "" {
			return c.Status(fiber.StatusBadRequest).JSON(NewHttpError[any](ErrNoEmailProvided, nil))
		}

		err := make(chan error)

		go subscription.Check(user.Email, user.UsedSpace, err)

		if err := <-err; err != nil {
			log.Error().Err(err).Msgf("access middleware: access denied: %v", err)
			return c.Status(fiber.StatusForbidden).JSON(NewHttpError[any](ErrAccessDenied, nil))
		}

		return c.Next()
	}
}
