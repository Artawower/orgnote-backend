package handlers

import (
	"moonbrain/app/models"
)

type CreatingNote struct {
	ID       string          `json:"id" form:"id"`
	Content  string          `json:"content" form:"content"`
	Meta     models.NoteMeta `json:"meta" form:"meta"`
	FilePath []string        `json:"filePath" form:"filePath"`
}

func mapCreatingNoteToNote(note CreatingNote) models.Note {
	return models.Note{
		ID:       note.ID,
		Content:  note.Content,
		Meta:     note.Meta,
		FilePath: note.FilePath,
	}
}

func mapCreatingNotesToNotes(notes []CreatingNote) (mappedNotes []models.Note) {
	mappedNotes = make([]models.Note, len(notes))
	for _, n := range notes {
		mappedNotes = append(mappedNotes, mapCreatingNoteToNote(n))
	}
	return
}
