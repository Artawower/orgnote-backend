package services

import "orgnote/app/models"

// TODO: master move to handler layer
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

func mapToPublicNote(note *models.Note, user *models.User, my bool) *models.PublicNote {
	u := mapToPublicUserInfo(user)
	return &models.PublicNote{
		ID:        note.ExternalID,
		Content:   note.Content,
		Meta:      note.Meta,
		FilePath:  note.FilePath,
		Author:    *u,
		UpdatedAt: note.UpdatedAt,
		IsMy:      my,
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
	}

}
