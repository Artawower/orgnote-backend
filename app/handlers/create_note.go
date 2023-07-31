package handlers

import (
	"moonbrain/app/models"
	"time"
)

type CreatingNote struct {
	ID        string          `json:"id" form:"id"`
	Content   string          `json:"content" form:"content"`
	Meta      models.NoteMeta `json:"meta" form:"meta"`
	FilePath  []string        `json:"filePath" form:"filePath"`
	UpdatedAt time.Time       `json:"updatedAt" form:"updatedAt"`
}

func mapCreatingNoteToNote(note CreatingNote) models.Note {
	return models.Note{
		ID:        note.ID,
		Content:   note.Content,
		Meta:      note.Meta,
		FilePath:  note.FilePath,
		UpdatedAt: note.UpdatedAt,
	}
}

func mapCreatingNotesToNotes(notes []CreatingNote) (mappedNotes []models.Note) {

	mappedNotes = make([]models.Note, 0, len(notes))
	for _, n := range notes {
		mappedNotes = append(mappedNotes, mapCreatingNoteToNote(n))
	}
	return
}
