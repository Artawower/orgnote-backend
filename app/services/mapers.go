package services

import "moonbrain/app/models"

func mapToPublicUserInfo(user *models.User) *models.PublicUser {
	return &models.PublicUser{
		ID:         user.ID.Hex(),
		Name:       user.Name,
		NickName:   user.NickName,
		AvatarURL:  user.AvatarURL,
		Email:      user.Email,
		ProfileURL: user.ProfileURL,
	}
}

func mapToPublicNote(note *models.Note, user *models.User) *models.PublicNote {
	u := mapToPublicUserInfo(user)
	return &models.PublicNote{
		ID:       note.ID,
		Content:  note.Content,
		Meta:     note.Meta,
		FilePath: note.FilePath,
		Author:   *u,
	}
}
