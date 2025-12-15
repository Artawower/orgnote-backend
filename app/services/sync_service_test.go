package services

import (
	"orgnote/app/infrastructure"
	"orgnote/app/models"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockBlobStorage struct {
	blobs map[string][]byte
}

func newMockBlobStorage() *mockBlobStorage {
	return &mockBlobStorage{blobs: make(map[string][]byte)}
}

func (m *mockBlobStorage) Upload(userID string, contentHash string, content []byte) error {
	m.blobs[userID+"/"+contentHash] = content
	return nil
}

func (m *mockBlobStorage) Download(userID string, contentHash string) ([]byte, error) {
	content, ok := m.blobs[userID+"/"+contentHash]
	if !ok {
		return nil, nil
	}
	return content, nil
}

func (m *mockBlobStorage) Exists(userID string, contentHash string) (bool, error) {
	_, ok := m.blobs[userID+"/"+contentHash]
	return ok, nil
}

func (m *mockBlobStorage) Delete(userID string, contentHash string) error {
	delete(m.blobs, userID+"/"+contentHash)
	return nil
}

func (m *mockBlobStorage) ListBlobs(userID string) ([]infrastructure.BlobInfo, error) {
	return nil, nil
}

type mockEventSender struct{}

func (m *mockEventSender) Emit(userID string, eventType string, payload interface{}, excludeSocketID string) {
}

type mockFileMetadataRepo struct {
	files     map[string]*models.FileMetadata
	totalSize int64
}

func newMockFileMetadataRepo() *mockFileMetadataRepo {
	return &mockFileMetadataRepo{files: make(map[string]*models.FileMetadata)}
}

func (m *mockFileMetadataRepo) Upsert(userID primitive.ObjectID, filePath string, contentHash string, fileSize int64, expectedVersion *int) (*models.FileMetadata, error) {
	metadata := &models.FileMetadata{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Path:        filePath,
		ContentHash: contentHash,
		Size:        fileSize,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Version:     1,
	}
	m.files[filePath] = metadata
	m.totalSize += fileSize
	return metadata, nil
}

func (m *mockFileMetadataRepo) GetTotalSize(userID primitive.ObjectID) (int64, error) {
	return m.totalSize, nil
}

func (m *mockFileMetadataRepo) HashExists(userID primitive.ObjectID, contentHash string) (bool, error) {
	for _, f := range m.files {
		if f.ContentHash == contentHash {
			return true, nil
		}
	}
	return false, nil
}

func TestUploadFile_Success(t *testing.T) {
	blobStorage := newMockBlobStorage()
	repo := newMockFileMetadataRepo()

	service := &SyncService{
		fileMetadataRepo:    nil,
		notificationService: NewNotificationService(&mockEventSender{}),
		blobStorage:         blobStorage,
		maxFileSize:         1024 * 1024,
		tombstoneTTL:        30 * 24 * time.Hour,
	}

	content := []byte("test content")
	hash := infrastructure.ComputeHash(content)

	err := blobStorage.Upload("user1", hash, content)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	downloaded, err := blobStorage.Download("user1", hash)
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}

	if string(downloaded) != string(content) {
		t.Errorf("content mismatch: got %s, want %s", downloaded, content)
	}

	_ = repo
	_ = service
}

func TestUploadFile_FileTooLarge(t *testing.T) {
	service := &SyncService{
		notificationService: NewNotificationService(&mockEventSender{}),
		maxFileSize:         10,
	}

	content := []byte("this content is too large")
	userID := primitive.NewObjectID()

	_, err := service.UploadFile(userID, "test.txt", content, "", 0, nil, "")

	if err != ErrFileTooLarge {
		t.Errorf("expected ErrFileTooLarge, got %v", err)
	}
}

func TestUploadFile_HashMismatch(t *testing.T) {
	service := &SyncService{
		notificationService: NewNotificationService(&mockEventSender{}),
		maxFileSize:         1024,
	}

	content := []byte("test")
	userID := primitive.NewObjectID()

	_, err := service.UploadFile(userID, "test.txt", content, "wrong-hash", 0, nil, "")

	if err != ErrHashMismatch {
		t.Errorf("expected ErrHashMismatch, got %v", err)
	}
}

func TestMapToChange_DeletedFile(t *testing.T) {
	service := &SyncService{}
	now := time.Now()

	metadata := models.FileMetadata{
		ID:        primitive.NewObjectID(),
		Path:      "deleted.txt",
		DeletedAt: &now,
		Version:   2,
	}

	change := service.mapToChange(metadata)

	if !change.Deleted {
		t.Error("expected Deleted to be true")
	}

	if change.ContentHash != nil {
		t.Error("expected ContentHash to be nil for deleted file")
	}
}

func TestMapToChange_ActiveFile(t *testing.T) {
	service := &SyncService{}

	metadata := models.FileMetadata{
		ID:          primitive.NewObjectID(),
		Path:        "active.txt",
		ContentHash: "abc123",
		DeletedAt:   nil,
		Version:     1,
	}

	change := service.mapToChange(metadata)

	if change.Deleted {
		t.Error("expected Deleted to be false")
	}

	if change.ContentHash == nil || *change.ContentHash != "abc123" {
		t.Error("expected ContentHash to be set")
	}
}
