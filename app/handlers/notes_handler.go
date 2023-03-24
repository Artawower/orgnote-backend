package handlers

import (
	"encoding/json"
	"moonbrain/app/models"
	"moonbrain/app/services"
	"net/http"

	_ "moonbrain/app/docs"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func collectNoteFromString(stringNote string) (models.Note, error) {
	note := models.Note{}
	err := json.Unmarshal([]byte(stringNote), &note)
	if err != nil {
		log.Error().Err(err).Msg("Error while unmarshalling note")
		return note, err
	}
	return note, nil
}

func collectNotesFromStrings(stringNotes []string) ([]models.Note, []string) {
	notes := []models.Note{}
	errors := []string{}
	for _, strNote := range stringNotes {
		note, err := collectNoteFromString(strNote)
		if err != nil {
			// TODO master: add user friendly error message
			errors = append(errors, err.Error())
			continue
		}
		notes = append(notes, note)
	}
	return notes, errors
}

type NoteHandlers struct {
	noteService *services.NoteService
}

// TODO: master wait when swago will support generics :(

type SuccessGetNotesResponse struct {
	Notes []models.Note `json:"notes"`
}

// GetNote godoc
// @Summary      Get note
// @Description  get note by id
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Note ID"
// @Success      200  {object}  HttpResponse[models.Note, any]
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes/{id}  [get]
func (h *NoteHandlers) GetNote(c *fiber.Ctx) error {
	noteID := c.Params("id")

	ctxUser := c.Locals("user")

	var userID string

	if ctxUser != nil {
		userID = ctxUser.(*models.User).ID.Hex()
	}

	notes, err := h.noteService.GetNote(noteID, userID)
	if err != nil {
		log.Info().Err(err).Msg("note handler: get note: get by id")
		return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Couldn't get note, something went wrong", nil))
	}
	if notes == nil {
		return c.Status(http.StatusNotFound).JSON(NewHttpReponse[any, any](nil, nil))
	}
	return c.Status(http.StatusOK).JSON(NewHttpReponse[*models.PublicNote, any](notes, nil))
}

// GetNote godoc
// @Summary      Get notes
// @Description  Get all notes with optional filter
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        userId       query  string  false  "User ID"
// @Param        searchText   query  string  false  "Search text"
// @Param        limit        query  int  true  "Limit for pagination"
// @Param        offset       query  int  true  "Offset for pagination"
// @Success      200  {object}  HttpResponse[[]models.Note, models.Pagination]
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes/  [get]
func (h *NoteHandlers) GetNotes(c *fiber.Ctx) error {
	defaultLimit := int64(10)
	defaultOffset := int64(0)

	filter := new(models.NoteFilter)

	if err := c.QueryParser(filter); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NewHttpError("Incorrect input query", err))
	}

	ctxUser := c.Locals("user")

	includePrivateNotes := filter.UserID != nil && ctxUser != nil && ctxUser.(*models.User).ID.Hex() == *filter.UserID

	if filter.Limit == nil {
		filter.Limit = &defaultLimit
	}

	if filter.Offset == nil {
		filter.Offset = &defaultOffset
	}

	paginatedNotes, err := h.noteService.GetNotes(includePrivateNotes, *filter)
	if err != nil {
		log.Info().Err(err).Msgf("note handler: get notes: get %v", err)
		return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Couldn't get notes, something went wrong", nil))
	}
	return c.Status(http.StatusOK).JSON(
		NewHttpReponse(paginatedNotes.Data, models.Pagination{
			Limit:  paginatedNotes.Limit,
			Offset: paginatedNotes.Offset,
			Total:  paginatedNotes.Total,
		}))
}

// CreateNote godoc
// @Summary      Create note
// @Description  Create note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        note       body  models.Note  true  "Note model"
// @Success      200  {object}  any
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes/  [post]
func (h *NoteHandlers) CreateNote(c *fiber.Ctx) error {
	note := new(models.Note)

	if err := c.BodyParser(note); err != nil {
		log.Info().Err(err).Msg("note handler: post note: parse body")
		return c.Status(fiber.StatusInternalServerError).JSON(NewHttpError("Can't parse body", err))
	}

	err := h.noteService.CreateNote(*note)

	if err != nil {
		log.Info().Err(err).Msgf("note handler: post note: create %v", err)
		return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Can't create note:(", nil))
	}
	return c.Status(http.StatusOK).JSON(nil)
}

// UpserNotes godoc
// @Summary      Upsert notes
// @Description  Bulk update or insert notes
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        notes body []models.Note  true  "Notes list"
// @Success      200  {object}  any
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes/bulk-upsert  [put]
func (h *NoteHandlers) UpsertNotes(c *fiber.Ctx) error {

	if form, err := c.MultipartForm(); err == nil {

		log.Info().Err(err).Msg("note handler: put notes: parse body")
		// files := form.File["files"]
		rawNotes, ok := form.Value["notes"]
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Notes doesn't provided", nil))
		}
		notes, errors := collectNotesFromStrings(rawNotes)
		if len(errors) > 0 {
			// TODO: master add errors exposing to real life.
			log.Error().Err(err).Msg("note handler: put notes: collect notes")
		}
		user := c.Locals("user").(*models.User)
		err = h.noteService.BulkCreateOrUpdate(user.ID.Hex(), notes)
		if err != nil {
			log.Warn().Msgf("note handlers: save notes: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Can't create notes", nil))
		}
		files := form.File["files"]
		log.Info().Msgf("notes: %v", files)

		err := h.noteService.UploadImages(files)
		if err != nil {
			// TODO: master error handling here
			return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Can't upload images", nil))
		}
		return c.Status(http.StatusOK).JSON(nil)
	}

	return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Can't parse multipart form data", nil))
}

// GetNoteGraph godoc
// @Summary      Get notes graph
// @Description  Return graph model with links between connected notes
// @Tags         notes
// @Accept       json
// @Produce      json
// @Success      200  {object}  handlers.HttpResponse[models.NoteGraph, any]
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes/graph  [get]
func (h *NoteHandlers) GetNoteGraph(c *fiber.Ctx) error {
	ctxUser := c.Locals("user")

	if ctxUser == nil {
		return c.Status(http.StatusNotFound).Send(nil)
	}

	graph, err := h.noteService.GetNoteGraph(ctxUser.(*models.User).ID.Hex())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Couldn't get note graph", nil))
	}

	return c.Status(http.StatusOK).JSON(NewHttpReponse[*models.NoteGraph, any](graph, nil))
}

func RegisterNoteHandler(app fiber.Router, noteService *services.NoteService, authMiddleware func(*fiber.Ctx) error) {
	noteHandlers := &NoteHandlers{
		noteService: noteService,
	}
	app.Get("/notes/graph", authMiddleware, noteHandlers.GetNoteGraph)
	app.Get("/notes/:id", noteHandlers.GetNote)
	app.Get("/notes", noteHandlers.GetNotes)
	app.Post("/notes", authMiddleware, noteHandlers.CreateNote)
	app.Put("/notes/bulk-upsert", authMiddleware, noteHandlers.UpsertNotes)
}
