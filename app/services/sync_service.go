package services

import (
	"errors"
	"fmt"
	"orgnote/app/infrastructure"
	"orgnote/app/models"
	"orgnote/app/repositories"
	"sync"
	"time"

	"orgnote/app/tools"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Syncer interface {
	GetChanges(userID primitive.ObjectID, since time.Time, limit int, cursor *string) (*ChangesResult, error)
	UploadFile(userID primitive.ObjectID, filePath string, content []byte, clientHash string, spaceLimit int64, expectedVersion *int) (*UploadResult, error)
	DownloadFile(userID primitive.ObjectID, filePath string) ([]byte, *models.FileMetadata, error)
	DeleteFile(userID primitive.ObjectID, filePath string, expectedVersion *int) (*models.FileMetadata, error)
	RunGarbageCollection(userID primitive.ObjectID) error
}

type SyncService struct {
	fileMetadataRepo *repositories.FileMetadataRepository
	blobStorage      infrastructure.BlobStorage
	maxFileSize      int64
	tombstoneTTL     time.Duration
}

type SyncServiceConfig struct {
	MaxFileSize  int64
	TombstoneTTL time.Duration
}

type ChangesResult struct {
	Changes    []models.FileChange
	Cursor     *string
	HasMore    bool
	ServerTime time.Time
}

type UploadResult struct {
	Metadata *models.FileMetadata
	Uploaded bool
}

func NewSyncService(
	fileMetadataRepo *repositories.FileMetadataRepository,
	blobStorage infrastructure.BlobStorage,
	config SyncServiceConfig,
) *SyncService {
	return &SyncService{
		fileMetadataRepo: fileMetadataRepo,
		blobStorage:      blobStorage,
		maxFileSize:      config.MaxFileSize,
		tombstoneTTL:     config.TombstoneTTL,
	}
}

func (s *SyncService) GetChanges(userID primitive.ObjectID, since time.Time, limit int, cursor *string) (*ChangesResult, error) {
	result, err := s.fileMetadataRepo.GetChanges(userID, since, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("sync service: get changes: %v", err)
	}

	return &ChangesResult{
		Changes:    s.mapToChanges(result.Files),
		Cursor:     result.NextCursor,
		HasMore:    result.HasMore,
		ServerTime: time.Now(),
	}, nil
}

func (s *SyncService) UploadFile(userID primitive.ObjectID, filePath string, content []byte, clientHash string, spaceLimit int64, expectedVersion *int) (*UploadResult, error) {
	filePath = tools.NormalizeFilePath(filePath)

	if int64(len(content)) > s.maxFileSize {
		return nil, ErrFileTooLarge
	}

	computedHash := infrastructure.ComputeHash(content)
	if clientHash != "" && clientHash != computedHash {
		return nil, ErrHashMismatch
	}

	currentUsage, err := s.fileMetadataRepo.GetTotalSize(userID)
	if err != nil {
		return nil, fmt.Errorf("sync service: upload: get total size: %v", err)
	}

	if spaceLimit > 0 && currentUsage+int64(len(content)) > spaceLimit {
		return nil, ErrStorageQuotaExceeded
	}

	uploaded, err := s.uploadBlobIfNeeded(userID.Hex(), computedHash, content)
	if err != nil {
		return nil, fmt.Errorf("sync service: upload: %v", err)
	}

	metadata, err := s.fileMetadataRepo.Upsert(userID, filePath, computedHash, int64(len(content)), expectedVersion)
	if versionErr, ok := err.(*repositories.VersionMismatchError); ok {
		return nil, &VersionMismatchError{
			Path:          versionErr.Path,
			ServerVersion: versionErr.ServerVersion,
		}
	}
	if errors.Is(err, repositories.ErrVersionMismatch) {
		return nil, ErrVersionMismatch
	}
	if err != nil {
		return nil, fmt.Errorf("sync service: upload: upsert metadata: %v", err)
	}

	return &UploadResult{
		Metadata: metadata,
		Uploaded: uploaded,
	}, nil
}

func (s *SyncService) uploadBlobIfNeeded(userID string, contentHash string, content []byte) (bool, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user id: %v", err)
	}

	exists, err := s.fileMetadataRepo.HashExists(userObjectID, contentHash)
	if err != nil {
		return false, fmt.Errorf("hash check: %v", err)
	}

	if exists {
		return false, nil
	}

	if err := s.blobStorage.Upload(userID, contentHash, content); err != nil {
		return false, fmt.Errorf("store blob: %v", err)
	}

	return true, nil
}

func (s *SyncService) mapToChanges(files []models.FileMetadata) []models.FileChange {
	changes := make([]models.FileChange, 0, len(files))
	for _, metadata := range files {
		changes = append(changes, s.mapToChange(metadata))
	}
	return changes
}

func (s *SyncService) mapToChange(metadata models.FileMetadata) models.FileChange {
	change := models.FileChange{
		ID:        metadata.ID.Hex(),
		Path:      metadata.Path,
		Size:      metadata.Size,
		UpdatedAt: metadata.UpdatedAt,
		Version:   metadata.Version,
		Deleted:   metadata.DeletedAt != nil,
		DeletedAt: metadata.DeletedAt,
	}

	if metadata.DeletedAt == nil {
		change.ContentHash = &metadata.ContentHash
	}

	return change
}

func (s *SyncService) DownloadFile(userID primitive.ObjectID, filePath string) ([]byte, *models.FileMetadata, error) {
	filePath = tools.NormalizeFilePath(filePath)

	metadata, err := s.fileMetadataRepo.GetByPath(userID, filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("sync service: download: get metadata: %v", err)
	}

	if metadata == nil {
		return nil, nil, nil
	}

	if metadata.DeletedAt != nil {
		return nil, nil, ErrFileDeleted
	}

	content, err := s.blobStorage.Download(userID.Hex(), metadata.ContentHash)
	if err != nil {
		return nil, nil, fmt.Errorf("sync service: download: get blob: %v", err)
	}

	if content == nil {
		return nil, nil, ErrBlobNotFound
	}

	return content, metadata, nil
}

func (s *SyncService) DeleteFile(userID primitive.ObjectID, filePath string, expectedVersion *int) (*models.FileMetadata, error) {
	filePath = tools.NormalizeFilePath(filePath)

	metadata, err := s.fileMetadataRepo.SoftDeleteByPath(userID, filePath, expectedVersion)
	if versionErr, ok := err.(*repositories.VersionMismatchError); ok {
		return nil, &VersionMismatchError{
			Path:          versionErr.Path,
			ServerVersion: versionErr.ServerVersion,
		}
	}
	if err != nil {
		return nil, fmt.Errorf("sync service: delete: %v", err)
	}

	return metadata, nil
}

func (s *SyncService) RunGarbageCollection(userID primitive.ObjectID) error {
	referencedHashes, err := s.fileMetadataRepo.GetReferencedHashes(userID)
	if err != nil {
		return fmt.Errorf("sync service: gc: get referenced hashes: %v", err)
	}

	referencedSet := make(map[string]bool)
	for _, hash := range referencedHashes {
		referencedSet[hash] = true
	}

	blobs, err := s.blobStorage.ListBlobs(userID.Hex())
	if err != nil {
		return fmt.Errorf("sync service: gc: list blobs: %v", err)
	}

	gcThreshold := time.Now().Add(-24 * time.Hour)
	s.deleteUnreferencedBlobs(userID.Hex(), blobs, referencedSet, gcThreshold)

	_, err = s.fileMetadataRepo.CleanOldTombstones(userID, s.tombstoneTTL)
	if err != nil {
		return fmt.Errorf("sync service: gc: clean tombstones: %v", err)
	}

	return nil
}

func (s *SyncService) deleteUnreferencedBlobs(
	userID string,
	blobs []infrastructure.BlobInfo,
	referencedSet map[string]bool,
	gcThreshold time.Time,
) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5)

	for _, blob := range blobs {
		if referencedSet[blob.Hash] {
			continue
		}
		if blob.CreatedAt.After(gcThreshold) {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(hash string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			if err := s.blobStorage.Delete(userID, hash); err != nil {
				log.Error().Err(err).Str("userId", userID).Str("hash", hash).Msg("sync service: gc: failed to delete blob")
			}
		}(blob.Hash)
	}

	wg.Wait()
}
