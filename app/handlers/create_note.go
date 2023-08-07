package handlers

import (
	"moonbrain/app/models"
	"time"

	"github.com/thoas/go-funk"
)

type CreatingNote struct {
	ID        string          `json:"id" form:"id"`
	Content   string          `json:"content" form:"content"`
	Meta      models.NoteMeta `json:"meta" form:"meta"`
	FilePath  []string        `json:"filePath" form:"filePath"`
	UpdatedAt time.Time       `json:"updatedAt" form:"updatedAt"`
	CreatedAt time.Time       `json:"createdAt" form:"createdAt"`
}

func mapCreatingNoteToNote(note CreatingNote) models.Note {
	return models.Note{
		ExternalID: note.ID,
		Content:    note.Content,
		Meta:       note.Meta,
		FilePath:   note.FilePath,
		UpdatedAt:  note.UpdatedAt,
		CreatedAt:  note.CreatedAt,
	}
}

func mapCreatingNotesToNotes(notes []CreatingNote) (mappedNotes []models.Note) {
	mappedNotes = (funk.Map(notes, func(n CreatingNote) models.Note {
		return mapCreatingNoteToNote(n)
	})).([]models.Note)
	return
}

func mapNoteToDeletedNote(note models.Note) DeletedNote {
	return DeletedNote{
		ID:       note.ExternalID,
		FilePath: note.FilePath,
	}
}
func mapNotesToDeletedNotes(notes []models.Note) (mappedNotes []DeletedNote) {
	return (funk.Map(notes, mapNoteToDeletedNote)).([]DeletedNote)
}
