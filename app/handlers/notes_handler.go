package handlers

import (
	"fmt"
	"moonbrain/app/models"
	"moonbrain/app/services"
	"net/http"
	"time"

	_ "moonbrain/app/docs"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
)

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
// @Success      200  {object}  HttpResponse[models.PublicNote, any]
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
		return c.Status(http.StatusNotFound).JSON(NewHttpResponse[any, any](nil, nil))
	}
	return c.Status(http.StatusOK).JSON(NewHttpResponse[*models.PublicNote, any](notes, nil))
}

// DeleteNotes godoc
// @Summary      Delete notes
// @Description  Mark notes as deleted by provided list of ids
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        ids   body     []string  true  "List of ids of deleted notes"
// @Success      200  {object}  HttpResponse[any, any]
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes [delete]
func (h *NoteHandlers) DeleteNotes(c *fiber.Ctx) error {
	notesIDs := []string{}
	err := c.BodyParser(&notesIDs)
	if err != nil {
		log.Info().Err(err).Msg("note handler: delete notes: body parser")
		return c.Status(http.StatusBadRequest).JSON(NewHttpError[any]("Couldn't parse body, something went wrong", nil))
	}
	h.noteService.DeleteNotes(notesIDs)
	return nil
}

type GetNotesFilter struct {
	Limit          *int64     `json:"limit" extensions:"x-order=1"`
	Offset         *int64     `json:"offset" extensions:"x-order=2"`
	UserID         *string    `json:"userId" extensions:"x-order=3"` // User id of which notes to load
	SearchText     *string    `json:"searchText" extensions:"x-order=4"`
	My             *bool      `json:"my" extensions:"x-order=5"` // Load all my own notes (user will be used from provided token)
	From           *time.Time `json:"from" extensions:"x-order=6"`
	IncludeDeleted *bool      `json:"includeDeleted" extensions:"x-order=7"`
}

var (
	defaultLimit  = int64(10)
	defaultOffset = int64(0)
)

func buildNotesFilter(user *models.User, filter *GetNotesFilter) *models.NoteFilter {
	if user != nil && filter.My != nil && *filter.My {
		userId := user.ID.Hex()
		filter.UserID = &userId
	}

	var published *bool
	if filter.UserID == nil || user != nil && *filter.UserID != user.ID.Hex() {
		pub := true
		published = &pub
	}

	if filter.Limit == nil {
		filter.Limit = &defaultLimit
	}

	if filter.Offset == nil {
		filter.Offset = &defaultOffset
	}

	return &models.NoteFilter{
		Limit:          filter.Limit,
		Offset:         filter.Offset,
		UserID:         filter.UserID,
		SearchText:     filter.SearchText,
		Published:      published,
		From:           filter.From,
		IncludeDeleted: filter.IncludeDeleted,
	}
}

// GetNote godoc
// @Summary      Get notes
// @Description  Get all notes with optional filter
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        filter       query  GetNotesFilter false "Filter"
// @Success      200  {object}  HttpResponse[[]models.PublicNote, models.Pagination]
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes/  [get]
func (h *NoteHandlers) GetNotes(c *fiber.Ctx) error {
	filter := new(GetNotesFilter)
	if err := c.QueryParser(filter); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NewHttpError("Incorrect input query", err))
	}

	ctxUser := c.Locals("user")

	serviceFilter := buildNotesFilter(ctxUser.(*models.User), filter)

	paginatedNotes, err := h.noteService.GetNotes(*serviceFilter)
	if err != nil {
		log.Info().Err(err).Msgf("note handler: get notes: get %v", err)
		return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Couldn't get notes, something went wrong", nil))
	}

	return c.Status(http.StatusOK).JSON(
		NewHttpResponse(paginatedNotes.Data, models.Pagination{
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
// @Param        note       body  CreatingNote  true  "Note model"
// @Success      200  {object}  any
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes/  [post]
func (h *NoteHandlers) CreateNote(c *fiber.Ctx) error {
	note := new(CreatingNote)

	if err := c.BodyParser(note); err != nil {
		log.Info().Msgf("note handler: post note: parse body: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(NewHttpError[any]("Can't parse body", nil))
	}

	author := c.Locals("user").(*models.User)
	n := mapCreatingNoteToNote(*note)
	n.AuthorID = author.ID.Hex()
	err := h.noteService.CreateNote(n)

	if err != nil {
		log.Info().Err(err).Msgf("note handler: post note: create %v", err)
		return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Can't create note", nil))
	}
	return c.Status(http.StatusOK).JSON(nil)
}

// UpserNotes godoc
// @Summary      Upsert notes
// @Description  Bulk update or insert notes
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        notes body []CreatingNote true "List of crated notes"
// @Success      200  {object}  any
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes/bulk-upsert  [put]
func (h *NoteHandlers) UpsertNotes(c *fiber.Ctx) error {
	notesForCreate := []CreatingNote{}

	if err := c.BodyParser(&notesForCreate); err != nil {
		log.Error().Err(err).Msgf("note handler: upsert notes: parse body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(NewHttpError[any]("Couldn't parse body, something went wrong", nil))
	}

	user := c.Locals("user").(*models.User)
	notes := mapCreatingNotesToNotes(notesForCreate)
	log.Info().Msgf("note handler: post note: create note id, note external id: %v", notes[0].ID, notes[0].ExternalID)

	err := h.noteService.BulkCreateOrUpdate(user.ID.Hex(), notes)
	if err != nil {
		log.Warn().Msgf("note handlers: save notes: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Can't create notes", nil))
	}
	return c.Status(http.StatusOK).JSON(nil)
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

	return c.Status(http.StatusOK).JSON(NewHttpResponse[*models.NoteGraph, any](graph, nil))
}

type SyncNotesRequest struct {
	Timestamp       time.Time      `json:"timestamp"`
	Notes           []CreatingNote `json:"notes"`
	DeletedNotesIDs []string       `json:"deletedNotesIds"`
}

type DeletedNote struct {
	ID       string   `json:"id"`
	FilePath []string `json:"filePath"`
}

type SyncNotesResponse struct {
	Notes        []models.PublicNote `json:"notes"`
	DeletedNotes []DeletedNote       `json:"deletedNotes"`
}

// SyncNotes godoc
// @Summary      Synchronize notes
// @Description  Synchronize notes with specific timestamp
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        data  body     SyncNotesRequest  true  "Sync notes request"
// @Success      200  {object}  HttpResponse[SyncNotesResponse, any]
// @Failure      400  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /notes/sync  [post]
func (h *NoteHandlers) SyncNotes(c *fiber.Ctx) error {
	ctxUser := c.Locals("user")
	userID := ctxUser.(*models.User).ID.Hex()

	// TODO: master validator!
	params := new(SyncNotesRequest)
	if err := c.BodyParser(params); err != nil {
		log.Info().Msgf("note handler: sync notes: parse body: %v", err)
		return fmt.Errorf("can't parse body")
	}
	notesToSync := mapCreatingNotesToNotes(params.Notes)
	notes, err := h.noteService.SyncNotes(notesToSync, params.DeletedNotesIDs, params.Timestamp, userID)

	if err != nil {
		log.Info().Err(err).Msg("note handler: sync notes")
		return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Couldn't sync notes", nil))
	}

	deletedNotes, err := h.noteService.GetDeletedNotes(userID, params.Timestamp)

	if err != nil {
		log.Info().Err(err).Msg("note handler: sync notes")
		return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any]("Couldn't sync notes", nil))
	}

	log.Info().Msgf("deleted notes: %v", funk.Map(deletedNotes, func(note models.Note) string {
		return note.ExternalID
	}))

	syncNotesResponse := SyncNotesResponse{
		Notes:        mapNotesToPublicNotes(notes, *ctxUser.(*models.User)),
		DeletedNotes: mapNotesToDeletedNotes(deletedNotes),
	}

	return c.Status(http.StatusOK).JSON(NewHttpResponse[SyncNotesResponse, any](syncNotesResponse, nil))
}

func RegisterNoteHandler(app fiber.Router, noteService *services.NoteService, authMiddleware func(*fiber.Ctx) error) {
	noteHandlers := &NoteHandlers{
		noteService: noteService,
	}
	app.Get("/notes/graph", authMiddleware, noteHandlers.GetNoteGraph)
	app.Get("/notes/:id", noteHandlers.GetNote)
	app.Get("/notes", noteHandlers.GetNotes)
	app.Post("/notes/sync", authMiddleware, noteHandlers.SyncNotes)
	app.Post("/notes", authMiddleware, noteHandlers.CreateNote)
	app.Put("/notes/bulk-upsert", authMiddleware, noteHandlers.UpsertNotes)
	app.Delete("/notes", authMiddleware, noteHandlers.DeleteNotes)
}
