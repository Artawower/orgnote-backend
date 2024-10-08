package services

import "orgnote/app/models"

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

func mapToPublicNote(note *models.Note, user *models.User, isMy bool) *models.PublicNote {
	u := mapToPublicUserInfo(user)
	return &models.PublicNote{
		ID:             note.ExternalID,
		Content:        note.Content,
		Meta:           note.Meta,
		FilePath:       note.FilePath,
		Author:         *u,
		UpdatedAt:      note.UpdatedAt,
		EncryptionType: note.EncryptionType,
		Encrypted:      note.Encrypted,
		TouchedAt:      note.TouchedAt,
		IsMy:           isMy,
	}
}

func mapToUserPersonalInfo(user *models.User) *models.UserPersonalInfo {
	return &models.UserPersonalInfo{
		ID:         user.ID.Hex(),
		Name:       user.Name,
		NickName:   user.NickName,
		AvatarURL:  user.AvatarURL,
		Email:      user.Email,
		ProfileURL: user.ProfileURL,
		SpaceLimit: user.SpaceLimit,
		UsedSpace:  user.UsedSpace,
		Active:     user.Active,
	}

}

func mapNotesToPublicNotes(notes []models.Note, user *models.User, isMy bool) (mappedNotes []models.PublicNote) {
	mappedNotes = make([]models.PublicNote, 0, len(notes))
	for _, n := range notes {
		mappedNotes = append(mappedNotes, *mapToPublicNote(&n, user, isMy))
	}
	return
}
