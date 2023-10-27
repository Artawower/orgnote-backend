package services

import (
	"fmt"
	"orgnote/app/models"
	"orgnote/app/repositories"
	"time"

	"github.com/rs/zerolog/log"
)

type NoteFileStorage interface {
	CalculateFileSize(folder string, fileName ...string) (int64, error)
}

type NoteService struct {
	noteRepository *repositories.NoteRepository
	userRepository *repositories.UserRepository
	tagRepository  *repositories.TagRepository
	fileStorage    NoteFileStorage
}

func NewNoteService(
	noteRepository *repositories.NoteRepository,
	userRepository *repositories.UserRepository,
	tagRepository *repositories.TagRepository,
	fileStorage NoteFileStorage,
) *NoteService {
	return &NoteService{
		noteRepository,
		userRepository,
		tagRepository,
		fileStorage,
	}
}

func (a *NoteService) CreateNote(note models.Note) error {
	err := a.noteRepository.AddNote(note)
	if err != nil {
		return err
	}
	return nil
}

func (n *NoteService) BulkCreateOrUpdate(userID string, notes []models.Note) error {
	defer func() {
		if len(notes) == 0 {
			return
		}
		go n.CalculateUserSpace(userID)
	}()

	filteredNotesWithID := []models.Note{}
	tags := []string{}
	for _, note := range notes {
		if note.ExternalID == "" {
			continue
		}
		note.AuthorID = userID
		filteredNotesWithID = append(filteredNotesWithID, models.Note{
			ID:         note.ID,
			ExternalID: note.ExternalID,
			AuthorID:   userID,
			Content:    note.Content,
			Meta:       note.Meta,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			FilePath:   note.FilePath,
			Views:      0,
			Likes:      0,
		})
		tags = append(tags, note.Meta.FileTags...)
	}
	// TODO: master add transaction here
	err := n.noteRepository.BulkUpsert(userID, filteredNotesWithID)
	if err != nil {
		return fmt.Errorf("note service: bulk create or update: could not bulk upsert notes: %v", err)
	}
	if len(tags) == 0 {
		return nil
	}
	err = n.tagRepository.BulkUpsert(tags)
	if err != nil {
		return fmt.Errorf("note service: bulk create or update: could not bulk upsert tags: %v", err)
	}

	return nil
}

// TODO: master return note. Move public note mapper into handlers.
// Repository should return full model with author (not an author id)
func (a *NoteService) GetNotes(filter models.NoteFilter, requestedUserId string) (*models.Paginated[models.PublicNote], error) {
	notes, err := a.noteRepository.GetNotes(filter)
	if err != nil {
		return nil, fmt.Errorf("note service: get notes: could not get notes: %v", err)
	}

	count, err := a.noteRepository.NotesCount(filter)
	if err != nil {
		return nil, fmt.Errorf("note service: upload images: get notes count: %v", err)
	}

	publicNotes := []models.PublicNote{}

	usersMap, err := a.getNotesUsers(notes)
	if err != nil {
		return nil, fmt.Errorf("note service: get notes: could not get users: %v", err)
	}

	for _, note := range notes {
		u := usersMap[note.AuthorID]
		my := note.AuthorID == requestedUserId
		publicNote := mapToPublicNote(&note, &u, my)
		publicNotes = append(publicNotes, *publicNote)
	}

	return &models.Paginated[models.PublicNote]{
		Limit:  *filter.Limit,
		Offset: *filter.Offset,
		Total:  count,
		Data:   publicNotes,
	}, nil
}

func (n *NoteService) GetDeletedNotes(userID string, deletedAt time.Time) ([]models.Note, error) {
	filter := models.NoteFilter{
		UserID:    &userID,
		DeletedAt: &deletedAt,
	}

	log.Info().Msgf("note service: get deleted notes: deleted at: %v", deletedAt)
	notes, err := n.noteRepository.GetNotes(filter)

	if err != nil {
		return nil, fmt.Errorf("note service: get deleted notes: could not get notes: %v", err)
	}

	return notes, nil
}

func (a *NoteService) getNotesUsers(notes []models.Note) (map[string]models.User, error) {
	if len(notes) == 0 {
		return map[string]models.User{}, nil
	}
	userIDSet := make(map[string]struct{})

	for _, note := range notes {
		userIDSet[note.AuthorID] = struct{}{}
	}

	userIDs := []string{}
	for k := range userIDSet {
		userIDs = append(userIDs, k)
	}

	users, err := a.userRepository.GetUsersByIDs(userIDs)

	if err != nil {
		return nil, fmt.Errorf("note service: get notes users: could not get users: %v", err)
	}

	usersMap := make(map[string]models.User)

	for _, u := range users {
		usersMap[u.ID.Hex()] = u
	}

	return usersMap, nil
}

func (a *NoteService) GetNote(id string, userID string) (*models.PublicNote, error) {
	note, err := a.noteRepository.GetNote(id, userID)
	if err != nil {
		return nil, fmt.Errorf("note service: get note: could not get note: %v", err)
	}
	if note == nil {
		return nil, nil
	}
	user, err := a.userRepository.GetByID(note.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("note service: get note: could not get user: %v", err)
	}
	myNote := userID == note.AuthorID
	publicNote := mapToPublicNote(note, user, myNote)
	return publicNote, nil
}

// TODO: master delete everything about graph. Redundant
func (n *NoteService) DeleteNotes(ids []string, authorID string) error {
	return n.noteRepository.MarkNotesAsDeleted(ids, authorID)
}

// TODO: master signature is too complex. Create a struct for params.
func (n *NoteService) SyncNotes(
	notes []models.Note,
	deletedNotesIDs []string,
	timestamp time.Time,
	authorID string,
) ([]models.Note, error) {
	filter := models.NoteFilter{
		From:           &timestamp,
		UserID:         &authorID,
		IncludeDeleted: new(bool),
	}

	err := n.noteRepository.DeleteOutdatedNotes(deletedNotesIDs, authorID, timestamp)

	if err != nil {
		return nil, fmt.Errorf("note service: sync notes: could not delete outdated notes: %v", err)
	}

	err = n.bulkUpdateOutdatedNotes(notes, authorID)

	if err != nil {
		return nil, err
	}

	notesFromLastSync, err := n.noteRepository.GetNotes(filter)

	if err != nil {
		return nil, fmt.Errorf("note service: sync notes: could not get notes: %v", err)
	}

	updatedNotes := n.excludeSameNotes(notesFromLastSync, notes)

	go n.CalculateUserSpace(authorID)
	return updatedNotes, nil
}

func (n *NoteService) bulkUpdateOutdatedNotes(notes []models.Note, authorID string) error {
	someNotesPresent := len(notes) > 0

	log.Info().Msgf("note service: notes length: %v", len(notes))

	if !someNotesPresent {
		return nil
	}
	err := n.noteRepository.BulkUpdateOutdated(notes, authorID)
	if err != nil {
		return fmt.Errorf("note service: sync notes: could not update outdated notes: %v", err)
	}
	return nil
}

func (n *NoteService) excludeSameNotes(srcNotes []models.Note, filterNotes []models.Note) []models.Note {
	filteredNotes := []models.Note{}

	for _, srcNote := range srcNotes {
		exists := false
		for _, filterNote := range filterNotes {
			if srcNote.ExternalID == filterNote.ExternalID && srcNote.UpdatedAt.Equal(filterNote.UpdatedAt) {
				exists = true
				break
			}
		}
		if !exists {
			filteredNotes = append(filteredNotes, srcNote)
		}
	}

	return filteredNotes
}

func (n *NoteService) CalculateUserSpace(userID string) error {
	spaceInfo, err := n.noteRepository.GetUsedSpaceInfo(userID)
	if err != nil {
		return fmt.Errorf("note service: calculate user space: could not calculate user space: %v", err)
	}
	log.Info().Msgf("note service: calculate user space: space info: %v", spaceInfo)

	usedFileSpace, err := n.fileStorage.CalculateFileSize(userID, spaceInfo.Files...)
	if err != nil {
		return fmt.Errorf("note service: calculate user space: could not calculate file size: %v", err)
	}

	totalUsedSpace := spaceInfo.UsedSpace + usedFileSpace

	err = n.userRepository.UpdateSpaceLimitInfo(userID, &totalUsedSpace, nil)

	if err != nil {
		return fmt.Errorf("note service: calculate user space: could not update used space: %v", err)
	}

	return nil
}
