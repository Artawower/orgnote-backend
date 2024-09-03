package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteHeading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}

type NoteLink struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type ConnectedNotes map[string]string

type category string

const (
	CategoryArticle  category = "article"
	CategoryBook     category = "book"
	CategorySchedule category = "schedule"
)

type NoteMeta struct {
	PreviewImg     *string         `json:"previewImg" bson:"previewImg"`
	Title          *string         `json:"title" bson:"title"`
	Description    *string         `json:"description" bson:"description"`
	Category       *category       `json:"category" bson:"category"`
	Headings       *[]NoteHeading  `json:"headings" bson:"headings"`
	ConnectedNotes *ConnectedNotes `json:"connectedNotes" bson:"connectedNotes"`
	Published      bool            `json:"published" bson:"published"`
	ExternalLinks  *[]NoteLink     `json:"externalLinks" bson:"externalLinks"`
	Startup        *string         `json:"startup" bson:"startup"`
	FileTags       []string        `json:"fileTags" bson:"fileTags"`
	Images         []string        `json:"images" bson:"images"`
}

type Note struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`                                                             // Generated ID for public notes
	EncryptionType *string            `json:"encryptionType" bson:"encryptionType" enums:"gpgKeys,gpgPassword,disabled"` // Encrypted note content
	Encrypted      bool               `json:"encrypted" bson:"encrypted"`
	ExternalID     string             `json:"externalId" bson:"externalId"` // Real note id. From source.
	AuthorID       string             `json:"authorId" bson:"authorId"`
	Content        string             `json:"content" bson:"content" binding:"required"`
	Meta           NoteMeta           `json:"meta" bson:"meta" binding:"required"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
	TouchedAt      time.Time          `json:"touchedAt" bson:"touchedAt"`
	LastSyncAt     time.Time          `json:"lastSyncAt" bson:"lastSyncAt"`
	FilePath       []string           `json:"filePath" bson:"filePath"`
	Views          int                `json:"views" bson:"views"`
	Likes          int                `json:"likes" bson:"likes"`
	DeletedAt      *time.Time         `json:"deletedAt" bson:"deletedAt"`
	Size           int64              `json:"size" bson"size"`
}

type PublicNote struct {
	ID             string     `json:"id"` // It's externalID from original note
	Author         PublicUser `json:"author" bson:"author"`
	EncryptionType *string    `json:"encryptionType" bson:"encryptionType" enums:"gpgKeys,gpgPassword,disabled"` // Encrypted note content
	Encrypted      bool       `json:"encrypted" bson:"encrypted"`
	Content        string     `json:"content" bson:"content" binding:"required"`
	Meta           NoteMeta   `json:"meta" binding:"required"`
	FilePath       []string   `json:"filePath"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	CreatedAt      time.Time  `json:"createdAt"`
	TouchedAt      time.Time  `json:"touchedAt"`
	IsMy           bool       `json:"isMy"`
	Size           int64      `json:"size" bson"size"`
}

type NoteFilter struct {
	Limit          *int64     `json:"limit"`
	Offset         *int64     `json:"offset"`
	UserID         *string    `json:"userId"`
	SearchText     *string    `json:"searchText"`
	Published      *bool      `json:"my"`
	From           *time.Time `json:"from" `
	IncludeDeleted *bool      `json:"includeDeleted"`
	DeletedAt      *time.Time `json:"deletedAt"`
}
