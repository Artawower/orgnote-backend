package handlers

import (
	"orgnote/app/models"

	"github.com/gofiber/fiber/v2"
)

type Subscription interface {
	Check(provider string, eternalID string, occupiedSpace int64, err chan<- error)
}

func NewAccessMiddleware(subscription Subscription) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)

		if user == (*models.User)(nil) {
			return c.Status(fiber.StatusBadRequest).JSON(NewHttpError[any](ErrAuthRequired, nil))
		}

		err := make(chan error)

		go subscription.Check(user.Provider, user.ExternalID, user.UsedSpace, err)

		if err := <-err; err != nil {
			return c.Status(fiber.StatusForbidden).JSON(NewHttpError[any](ErrAccessDenied, err.Error()))
		}

		return c.Next()
	}
}
