package services

import (
	"fmt"
	"moonbrain/app/models"
	"moonbrain/app/repositories"
	"moonbrain/app/tools"
	"time"
)

type NoteService struct {
	noteRepository *repositories.NoteRepository
	userRepository *repositories.UserRepository
	tagRepository  *repositories.TagRepository
}

func NewNoteService(
	noteRepository *repositories.NoteRepository,
	userRepository *repositories.UserRepository,
	tagRepository *repositories.TagRepository,
) *NoteService {
	return &NoteService{
		noteRepository: noteRepository,
		tagRepository:  tagRepository,
		userRepository: userRepository,
	}
}

func (a *NoteService) CreateNote(note models.Note) error {
	err := a.noteRepository.AddNote(note)
	if err != nil {
		return err
	}
	return nil
}

func (a *NoteService) BulkCreateOrUpdate(userID string, notes []models.Note) error {
	filteredNotesWithID := []models.Note{}
	tags := []string{}
	for _, note := range notes {
		if note.ID == "" {
			continue
		}
		note.AuthorID = userID
		filteredNotesWithID = append(filteredNotesWithID, models.Note{
			ID:        note.ID,
			AuthorID:  userID,
			Content:   note.Content,
			Meta:      note.Meta,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FilePath:  note.FilePath,
			Views:     0,
			Likes:     0,
		})
		tags = append(tags, note.Meta.FileTags...)
		go a.udpateNoteGraph(userID, note)
	}
	// TODO: master add transaction here
	err := a.noteRepository.BulkUpsert(userID, filteredNotesWithID)
	if err != nil {
		return fmt.Errorf("note service: bulk create or update: could not bulk upsert notes: %v", err)
	}
	if len(tags) == 0 {
		return nil
	}
	err = a.tagRepository.BulkUpsert(tags)
	if err != nil {
		return fmt.Errorf("note service: bulk create or update: could not bulk upsert tags: %v", err)
	}

	return nil
}

func (a *NoteService) GetNotes(includePrivate bool, filter models.NoteFilter) (*models.Paginated[models.PublicNote], error) {
	notes, err := a.noteRepository.GetNotes(includePrivate, filter)
	if err != nil {
		return nil, fmt.Errorf("note service: get notes: could not get notes: %v", err)
	}

	count, err := a.noteRepository.NotesCount(includePrivate, filter)
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
		publicNote := mapToPublicNote(&note, &u)
		publicNotes = append(publicNotes, *publicNote)
	}

	return &models.Paginated[models.PublicNote]{
		Limit:  *filter.Limit,
		Offset: *filter.Offset,
		Total:  count,
		Data:   publicNotes,
	}, nil
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
	publicNote := mapToPublicNote(note, user)
	return publicNote, nil
}

func (a *NoteService) GetNoteGraph(userID string) (*models.NoteGraph, error) {
	graph, err := a.userRepository.GetNoteGraph(userID)
	if err != nil {
		return nil, fmt.Errorf("note service: get note graph: could not get note graph: %v", err)
	}
	return graph, nil
}

func (a *NoteService) udpateNoteGraph(userID string, note models.Note) error {

	currentNoteNode := a.getGraphNoteNode(note)
	relatedLinks := a.getRelatedLinks(note)

	graphNoteLinks := repositories.GraphNoteLinks{
		Node:  currentNoteNode,
		Links: relatedLinks,
	}
	err := a.userRepository.UpsertGraphNode(userID, graphNoteLinks)
	if err != nil {
		// TODO: add this job to queue and log error
		return fmt.Errorf("note service: update note graph: upser graph node: %v", err)
	}
	return nil
}

func (a *NoteService) getGraphNoteNode(note models.Note) models.GraphNoteNode {
	weight := 0
	if note.Meta.LinkedArticles != nil {
		weight = len(*note.Meta.LinkedArticles)
	}

	title := ""
	if note.Meta.Title != nil {
		title = *note.Meta.Title
	}

	return models.GraphNoteNode{
		ID:     note.ID,
		Title:  title,
		Weight: weight,
	}
}

func (a *NoteService) getRelatedLinks(note models.Note) (graphNoteLinks []models.GraphNoteLink) {
	graphNoteLinks = []models.GraphNoteLink{}
	if note.Meta.ExternalLinks == nil {
		return
	}
	for _, relation := range *note.Meta.LinkedArticles {

		realID, ok := tools.ExportLinkID(relation.Url)
		if !ok {
			continue
		}
		graphNoteLinks = append(graphNoteLinks, models.GraphNoteLink{
			Source: note.ID,
			Target: realID,
		})
	}

	return
}

func (n *NoteService) DeleteNotes(ids []string) error {
	return n.noteRepository.MarkNotesAsDeleted(ids)
}

func (n *NoteService) SyncNotes(notes []models.Note, timestamp time.Time, authorID string) ([]models.Note, error) {
	my := true
	filter := models.NoteFilter{
		My:     &my,
		From:   &timestamp,
		UserID: &authorID,
	}

	err := n.bulkUpdateOutdatedNotes(notes, authorID)

	if err != nil {
		return nil, err
	}

	notesFromLastSync, err := n.noteRepository.GetNotes(true, filter)
	if err != nil {
		return nil, fmt.Errorf("note service: sync notes: could not get notes: %v", err)
	}

	updatedNotes := n.excludeSameNotes(notesFromLastSync, notes)

	return updatedNotes, nil
}

func (n *NoteService) bulkUpdateOutdatedNotes(notes []models.Note, authorID string) error {
	someNotesPresent := len(notes) > 0

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
			if srcNote.ID == filterNote.ID && srcNote.UpdatedAt.Equal(filterNote.UpdatedAt) {
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
