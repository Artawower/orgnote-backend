package handlers

import "orgnote/app/models"

func mapToPublicUserInfo(user models.User) models.PublicUser {
	return models.PublicUser{
		ID:         user.ID.Hex(),
		Name:       user.Name,
		NickName:   user.NickName,
		AvatarURL:  user.AvatarURL,
		Email:      user.Email,
		ProfileURL: user.ProfileURL,
	}
}

func mapNoteToPublicNote(note models.Note, user models.User) models.PublicNote {
	u := mapToPublicUserInfo(user)
	return models.PublicNote{
		ID:        note.ExternalID,
		Content:   note.Content,
		Meta:      note.Meta,
		FilePath:  note.FilePath,
		Author:    u,
		UpdatedAt: note.UpdatedAt,
		IsMy:      user.ID.Hex() == note.AuthorID,
	}
}

func mapNotesToPublicNotes(notes []models.Note, user models.User) (mappedNotes []models.PublicNote) {
	mappedNotes = make([]models.PublicNote, 0, len(notes))
	for _, n := range notes {
		mappedNotes = append(mappedNotes, mapNoteToPublicNote(n, user))
	}
	return
}
