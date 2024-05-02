package handlers

import (
	"orgnote/app/models"

	"github.com/gofiber/fiber/v2"
)

type OrgNoteMetaService interface {
	GetChangesFrom(version string) *models.OrgNoteClientUpdateInfo
}

type SystemInfoHandler struct {
	metaService OrgNoteMetaService
}

func NewSystemInfoHandler(metaService OrgNoteMetaService) *SystemInfoHandler {
	return &SystemInfoHandler{metaService: metaService}
}

// ClientUpdateInfo godoc
// @Summary      GetUpdatesFromVersion
// @Description
// @Tags         system info
// @Accept       json
// @Produce      json
// @Param        version path string true "provider"
// @Success      200  {object}  models.OrgNoteClientUpdateInfo
// @Failure      400  {object}  handlers.HttpError[any]
// @Failure      404  {object}  handlers.HttpError[any]
// @Failure      500  {object}  handlers.HttpError[any]
// @Router       /system-info/client-update/{version} [get]
func (s *SystemInfoHandler) LoginCallback(c *fiber.Ctx) error {
	version := c.Params("version")
	updateInfo := s.metaService.GetChangesFrom(version)
	if updateInfo == nil {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
	}
	return c.JSON(updateInfo)
}

func RegisterSystemInfoHandler(app fiber.Router, metaService OrgNoteMetaService) {
	handler := NewSystemInfoHandler(metaService)
	app.Get("/system-info/client-update/:version", handler.LoginCallback)
}
