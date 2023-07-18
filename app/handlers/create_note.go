package handlers

import (
	"moonbrain/app/models"
)

type CreatedNote struct {
	ID       string          `json:"id" form:"id"`
	Content  string          `json:"content" form:"content"`
	Meta     models.NoteMeta `json:"meta" form:"meta"`
	FilePath []string        `json:"filePath" form:"filePath"`
}

type UpsertNoteFormData struct {
	Notes CreatedNote `form:"notes" json:"notes"`
	Files string      `form:"files" json:"files" format:"binary"`
}

func mapCreatingNoteToNote(note CreatedNote) models.Note {
	return models.Note{
		ID:       note.ID,
		Content:  note.Content,
		Meta:     note.Meta,
		FilePath: note.FilePath,
	}
}

func mapCreatingNotesToNotes(notes []CreatedNote) (mappedNotes []models.Note) {
	mappedNotes = make([]models.Note, len(notes))
	for _, n := range notes {
		mappedNotes = append(mappedNotes, mapCreatingNoteToNote(n))
	}
	return
}
