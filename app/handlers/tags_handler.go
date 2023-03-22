package handlers

import (
	"moonbrain/app/services"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

 
// GetTags godoc
// @Summary      Get tags
// @Description  Return list of al registered tags
// @Tags         tags
// @Accept       json
// @Produce      json
// @Success      200  {object}  handlers.HttpResponse[[]string, any]
// @Failure      400  {object}  handlers.HttpError[any]
// @Failure      404  {object}  handlers.HttpError[any]
// @Failure      500  {object}  handlers.HttpError[any]
// @Router       /notes/graph  [get]
func RegisterTagHandler(app fiber.Router, tagService *services.TagService) {
	app.Get("/tags", func(c *fiber.Ctx) error {
		tags, err := tagService.GetTags()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any](err.Error(), nil))
		}
		return c.Status(http.StatusOK).JSON(NewHttpReponse[[]string, any](tags, nil))
	})
}
