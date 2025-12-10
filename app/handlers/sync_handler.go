package handlers

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"orgnote/app/models"
	"orgnote/app/services"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var reservedNamesPattern = regexp.MustCompile(`(?i)^(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(\.|$)`)

type SyncHandler struct {
	syncService *services.SyncService
	validate    *Validator
}

type GetChangesRequest struct {
	Since  int64  `query:"since" validate:"min=0"`
	Limit  int    `query:"limit" validate:"min=1,max=500"`
	Cursor string `query:"cursor"`
}

type UploadFileRequest struct {
	FilePath        string `validate:"required,min=1"`
	FileContent     []byte `validate:"required,min=1"`
	ClientHash      string
	ExpectedVersion *int
}

type FilePathParams struct {
	ID string `params:"id"`
}

type DeleteFileRequest struct {
	FilePathParams
	Version int `query:"version"`
}

func (h *SyncHandler) parseGetChangesRequest(c *fiber.Ctx) (*GetChangesRequest, error) {
	req := &GetChangesRequest{
		Limit: 100,
	}
	if err := c.QueryParser(req); err != nil {
		return nil, err
	}

	if errs := h.validate.Validate(req); len(errs) > 0 {
		return nil, fmt.Errorf("validation failed: %v", errs)
	}

	return req, nil
}

func (h *SyncHandler) parseUploadFileRequest(c *fiber.Ctx) (*UploadFileRequest, error) {
	filePath := c.FormValue("filePath")
	if err := validateFilePath(filePath); err != nil {
		return nil, err
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("file is required")
	}

	content, err := h.readFileContent(fileHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var expectedVersion *int
	if versionStr := c.FormValue("expectedVersion"); versionStr != "" {
		v, err := strconv.Atoi(versionStr)
		if err == nil && v > 0 {
			expectedVersion = &v
		}
	}

	req := &UploadFileRequest{
		FilePath:        filePath,
		FileContent:     content,
		ClientHash:      c.Get("X-Content-Hash"),
		ExpectedVersion: expectedVersion,
	}

	if errs := h.validate.Validate(req); len(errs) > 0 {
		return nil, fmt.Errorf("validation failed: %v", errs)
	}

	return req, nil
}

func validateFilePath(path string) error {
	if path == "" {
		return fmt.Errorf("filePath is required")
	}

	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("file path cannot be absolute")
	}

	if strings.Contains(path, "..") {
		return fmt.Errorf("file path cannot contain '..'")
	}

	for _, part := range strings.Split(path, "/") {
		if part == "" {
			return fmt.Errorf("file path contains empty segment")
		}
		if reservedNamesPattern.MatchString(part) {
			return fmt.Errorf("file path contains reserved name: %s", part)
		}
	}

	return nil
}

func (h *SyncHandler) readFileContent(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func (h *SyncHandler) parseFilePathParams(c *fiber.Ctx) (*FilePathParams, error) {
	req := &FilePathParams{}
	if err := c.ParamsParser(req); err != nil {
		return nil, err
	}
	if req.ID == "" {
		return nil, fmt.Errorf("file ID is required")
	}
	return req, nil
}

func (h *SyncHandler) parseDeleteFileRequest(c *fiber.Ctx) (*DeleteFileRequest, error) {
	req := &DeleteFileRequest{}
	if err := c.ParamsParser(req); err != nil {
		return nil, err
	}
	if err := c.QueryParser(req); err != nil {
		return nil, err
	}
	if req.ID == "" {
		return nil, fmt.Errorf("file ID is required")
	}
	return req, nil
}

func (h *SyncHandler) badRequest(c *fiber.Ctx, msg string) error {
	return c.Status(http.StatusBadRequest).JSON(NewHttpError[any](msg, nil))
}

func (h *SyncHandler) notFound(c *fiber.Ctx, msg string) error {
	return c.Status(http.StatusNotFound).JSON(NewHttpError[any](msg, nil))
}

func (h *SyncHandler) serverError(c *fiber.Ctx, err error, context string) error {
	log.Error().Err(err).Msg(context)
	return c.Status(http.StatusInternalServerError).JSON(NewHttpError[any](err.Error(), nil))
}

// GetChanges godoc
// @Summary      Get file changes
// @Description  Returns file changes since the specified timestamp
// @Tags         sync
// @Accept       json
// @Produce      json
// @Param        since query string false "ISO8601 timestamp for incremental sync"
// @Param        limit query int false "Maximum number of changes to return (default: 100, max: 500)"
// @Param        cursor query string false "Pagination cursor"
// @Success      200  {object}  HttpResponse[models.SyncChangesResponse, any]
// @Failure      400  {object}  HttpError[any]
// @Failure      401  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /sync/changes [get]
func (h *SyncHandler) GetChanges(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	req, err := h.parseGetChangesRequest(c)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	since := time.UnixMilli(req.Since)

	var cursor *string
	if req.Cursor != "" {
		cursor = &req.Cursor
	}

	resp, err := h.syncService.GetChanges(user.ID, since, req.Limit, cursor)
	if err != nil {
		return h.serverError(c, err, "sync handler: get changes")
	}

	response := &models.SyncChangesResponse{
		Changes:    resp.Changes,
		Cursor:     resp.Cursor,
		HasMore:    resp.HasMore,
		ServerTime: resp.ServerTime,
	}

	return c.JSON(NewHttpResponse[*models.SyncChangesResponse, any](response, nil))
}

// UploadFile godoc
// @Summary      Upload a file
// @Description  Upload a file to sync storage with content-addressable deduplication
// @Tags         sync
// @Accept       multipart/form-data
// @Produce      json
// @Param        filePath formData string true "Relative file path"
// @Param        file formData file true "File content"
// @Param        expectedVersion formData int false "Expected version for optimistic locking"
// @Param        X-Content-Hash header string false "SHA-256 hash for verification"
// @Success      200  {object}  HttpResponse[models.FileUploadResponse, any]
// @Failure      400  {object}  HttpError[any]
// @Failure      401  {object}  HttpError[any]
// @Failure      409  {object}  HttpError[any]
// @Failure      413  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /sync/files [put]
func (h *SyncHandler) UploadFile(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	req, err := h.parseUploadFileRequest(c)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	result, err := h.syncService.UploadFile(user.ID, req.FilePath, req.FileContent, req.ClientHash, user.SpaceLimit, req.ExpectedVersion)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrStorageQuotaExceeded):
			return c.Status(http.StatusRequestEntityTooLarge).JSON(NewHttpError[any]("storage limit exceeded", nil))
		case errors.Is(err, services.ErrFileTooLarge):
			return c.Status(http.StatusRequestEntityTooLarge).JSON(NewHttpError[any]("file too large", nil))
		case errors.Is(err, services.ErrHashMismatch):
			return h.badRequest(c, "hash mismatch")
		case errors.Is(err, services.ErrVersionMismatch):
			return c.Status(http.StatusConflict).JSON(NewHttpError[any]("version mismatch: file was modified", nil))
		default:
			return h.serverError(c, err, "sync handler: upload")
		}
	}

	response := &models.FileUploadResponse{
		ID:          result.Metadata.ID.Hex(),
		FilePath:    result.Metadata.FilePath,
		ContentHash: result.Metadata.ContentHash,
		FileSize:    result.Metadata.FileSize,
		UpdatedAt:   result.Metadata.UpdatedAt,
		Version:     result.Metadata.Version,
		Uploaded:    result.Uploaded,
	}

	return c.JSON(NewHttpResponse[*models.FileUploadResponse, any](response, nil))
}

// DownloadFile godoc
// @Summary      Download a file
// @Description  Download file content by file ID
// @Tags         sync
// @Produce      octet-stream
// @Param        id path string true "File ID"
// @Success      200  {file}  binary
// @Failure      400  {object}  HttpError[any]
// @Failure      401  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /sync/files/{id} [get]
func (h *SyncHandler) DownloadFile(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	params, err := h.parseFilePathParams(c)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	fileID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return h.badRequest(c, "invalid file ID")
	}

	content, metadata, err := h.syncService.DownloadFile(user.ID, fileID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFileDeleted):
			return h.notFound(c, "file deleted")
		case errors.Is(err, services.ErrBlobNotFound):
			return h.serverError(c, err, "sync handler: download: blob not found")
		default:
			return h.serverError(c, err, "sync handler: download")
		}
	}

	if metadata == nil {
		return h.notFound(c, "file not found")
	}

	c.Set("Content-Type", "application/octet-stream")
	c.Set("X-Content-Hash", metadata.ContentHash)
	c.Set("X-File-Path", metadata.FilePath)
	c.Set("Content-Length", strconv.FormatInt(metadata.FileSize, 10))

	return c.Send(content)
}

// DeleteFile godoc
// @Summary      Delete a file
// @Description  Soft delete a file (creates tombstone for sync)
// @Tags         sync
// @Accept       json
// @Produce      json
// @Param        id path string true "File ID"
// @Param        version query int false "Expected version for optimistic locking"
// @Success      200  {object}  HttpResponse[models.FileMetadata, any]
// @Failure      400  {object}  HttpError[any]
// @Failure      401  {object}  HttpError[any]
// @Failure      404  {object}  HttpError[any]
// @Failure      409  {object}  HttpError[any]
// @Failure      500  {object}  HttpError[any]
// @Router       /sync/files/{id} [delete]
func (h *SyncHandler) DeleteFile(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	req, err := h.parseDeleteFileRequest(c)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	fileID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return h.badRequest(c, "invalid file ID")
	}

	var expectedVersion *int
	if req.Version > 0 {
		expectedVersion = &req.Version
	}

	metadata, err := h.syncService.DeleteFile(user.ID, fileID, expectedVersion)
	if err != nil {
		return h.serverError(c, err, "sync handler: delete")
	}

	if metadata == nil {
		return h.notFound(c, "file not found or version mismatch")
	}

	return c.JSON(NewHttpResponse[*models.FileMetadata, any](metadata, nil))
}

func RegisterSyncHandler(
	app fiber.Router,
	syncService *services.SyncService,
	authMiddleware func(*fiber.Ctx) error,
	accessMiddleware func(*fiber.Ctx) error,
) {
	handler := &SyncHandler{
		syncService: syncService,
		validate:    NewValidator(),
	}

	app.Get("/sync/changes", authMiddleware, accessMiddleware, handler.GetChanges)
	app.Put("/sync/files", authMiddleware, accessMiddleware, handler.UploadFile)
	app.Get("/sync/files/:id", authMiddleware, accessMiddleware, handler.DownloadFile)
	app.Delete("/sync/files/:id", authMiddleware, accessMiddleware, handler.DeleteFile)
}
