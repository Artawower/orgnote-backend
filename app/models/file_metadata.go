package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileMetadata struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	FilePath    string             `json:"filePath" bson:"filePath"`
	ContentHash string             `json:"contentHash" bson:"contentHash"`
	FileSize    int64              `json:"fileSize" bson:"fileSize"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt   *time.Time         `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
	Version     int                `json:"version" bson:"version"`
}

// FileChange represents a single file change for sync
// @Description File change information for synchronization
type FileChange struct {
	ID          string     `json:"id" binding:"required" example:"507f1f77bcf86cd799439011"`
	FilePath    string     `json:"filePath" binding:"required" example:"notes/todo.org"`
	ContentHash *string    `json:"contentHash,omitempty" example:"a1b2c3d4e5f6..."`
	FileSize    int64      `json:"fileSize,omitempty" example:"1024"`
	UpdatedAt   time.Time  `json:"updatedAt" binding:"required" example:"2024-01-01T00:00:00Z"`
	Deleted     bool       `json:"deleted" binding:"required" example:"false"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	Version     int        `json:"version" binding:"required" example:"1"`
} // @name FileChange

// SyncChangesResponse represents the response for sync changes endpoint
// @Description Response containing file changes since last sync
type SyncChangesResponse struct {
	Changes    []FileChange `json:"changes" binding:"required"`
	Cursor     *string      `json:"cursor,omitempty" example:"507f1f77bcf86cd799439011"`
	HasMore    bool         `json:"hasMore" binding:"required" example:"false"`
	ServerTime time.Time    `json:"serverTime" binding:"required" example:"2024-01-01T00:00:00Z"`
} // @name SyncChangesResponse

// FileUploadResponse represents the response after file upload
// @Description Response after successful file upload
type FileUploadResponse struct {
	ID          string    `json:"id" binding:"required" example:"507f1f77bcf86cd799439011"`
	FilePath    string    `json:"filePath" binding:"required" example:"notes/todo.org"`
	ContentHash string    `json:"contentHash" binding:"required" example:"a1b2c3d4e5f6..."`
	FileSize    int64     `json:"fileSize" binding:"required" example:"1024"`
	UpdatedAt   time.Time `json:"updatedAt" binding:"required" example:"2024-01-01T00:00:00Z"`
	Version     int       `json:"version" binding:"required" example:"1"`
	Uploaded    bool      `json:"uploaded" binding:"required" example:"true"`
} // @name FileUploadResponse

type SyncCursor struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId" bson:"userId"`
	DeviceID   string             `json:"deviceId" bson:"deviceId"`
	LastSyncAt time.Time          `json:"lastSyncAt" bson:"lastSyncAt"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}
