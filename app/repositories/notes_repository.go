package repositories

import (
	"context"
	"errors"
	"fmt"
	"orgnote/app/models"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoteRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewNoteRepository(db *mongo.Database) *NoteRepository {
	noteRepo := &NoteRepository{db: db, collection: db.Collection("notes")}
	noteRepo.initIndexes()
	return noteRepo
}

func (a *NoteRepository) initIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := a.collection.Indexes().DropAll(ctx)
	if err != nil {
		log.Error().Msgf("note repository: failed to drop indexes: %v", err)
	}
	model := []mongo.IndexModel{
		{Keys: bson.D{
			bson.E{Key: "meta.title", Value: "text"},
			bson.E{Key: "meta.description", Value: "text"},
			bson.E{Key: "meta.tags", Value: "text"},
		}}}

	name, err := a.collection.Indexes().CreateMany(context.TODO(), model)
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("note repository: created indexes: %v", name)
}

func (a *NoteRepository) GetNotes(f models.NoteFilter) ([]models.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	notes := []models.Note{}
	filter := getNotesFilter(f)

	log.Info().Msgf("note repository: get notes: filter: %v", filter)

	findOptions := options.FindOptions{}

	if f.Limit != nil {
		findOptions.SetLimit(*f.Limit)
	}

	if f.Offset != nil {
		findOptions.SetSkip(*f.Offset)
	}

	findOptions.SetSort(bson.D{bson.E{Key: "createdAt", Value: -1}})

	cur, err := a.collection.Find(ctx, filter, &findOptions)
	if err != nil {
		return nil, fmt.Errorf("note repository: failed to get notes: %v", err)
	}

	for cur.Next(ctx) {
		var note models.Note
		err := cur.Decode(&note)
		if err != nil {
			return nil, fmt.Errorf("note repository: failed to decode note: %v", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (a *NoteRepository) NotesCount(f models.NoteFilter) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := getNotesFilter(f)
	count, err := a.collection.CountDocuments(ctx, filters)
	if err != nil {
		return 0, fmt.Errorf("note repository: failed to get notes count: %v", err)
	}
	return count, nil
}

func (a *NoteRepository) AddNote(note models.Note) error {
	// TODO: create note
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := a.collection.InsertOne(ctx, note)

	if err != nil {
		return fmt.Errorf("note repository: failed to add note: %v", err)
	}

	return nil
}

func (a *NoteRepository) BulkUpsert(userID string, notes []models.Note) error {
	if (len(notes)) == 0 {
		return errors.New("note repository: no notes to upsert")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notesModels := make([]mongo.WriteModel, len(notes))

	for i, note := range notes {
		// TODO: master id should be unique for each user
		notesModels[i] = mongo.NewUpdateOneModel().
			SetFilter(bson.M{"externalId": note.ExternalID, "authorId": userID}).
			SetUpdate(bson.M{
				"$set":         a.getUpdateNote(note),
				"$setOnInsert": bson.M{"_id": primitive.NewObjectID(), "createdAt": note.CreatedAt},
			}).
			SetUpsert(true)
	}

	_, err := a.collection.BulkWrite(ctx, notesModels)
	if err != nil {
		return fmt.Errorf("note repository: failed to bulk upsert notes: %v", err)
	}
	return nil
}

func (a *NoteRepository) getUpdateNote(note models.Note) bson.M {
	update := bson.M{
		"externalId": note.ExternalID,
		"authorId":   note.AuthorID,
		"content":    note.Content,
		"meta":       note.Meta,
		"updatedAt":  note.UpdatedAt,
		"views":      note.Views,
		"likes":      note.Likes,
		// "deletedAt": nil,
		"filePath": note.FilePath,
	}

	return update
}

func (a *NoteRepository) GetNote(externalID string, authorID string) (*models.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := a.collection.FindOne(
		ctx,
		bson.M{
			"externalId": externalID,
			"$or": bson.A{
				bson.M{"authorId": authorID},
				bson.M{"meta.published": true}}})

	err := res.Err()

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("note repository: failed to get note: %v", err)
	}

	var note models.Note
	err = res.Decode(&note)
	if err != nil {
		return nil, fmt.Errorf("note repository: failed to decode note: %v", err)
	}
	return &note, nil
}

func (n *NoteRepository) MarkNotesAsDeleted(noteIds []string, authorId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notesModel := make([]mongo.WriteModel, len(noteIds))

	for i, noteId := range noteIds {
		notesModel[i] = mongo.NewUpdateOneModel().
			SetFilter(bson.M{"externalId": noteId, "authorId": authorId}).
			SetUpdate(bson.M{"$set": bson.M{"deletedAt": time.Now()}})
	}

	_, err := n.collection.BulkWrite(ctx, notesModel)
	if err != nil {
		return fmt.Errorf("note repository: failed to bulk update notes: %v", err)
	}

	return nil
}

func (n *NoteRepository) BulkUpdateOutdated(notes []models.Note, authorID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notesModel := []mongo.WriteModel{}

	for _, note := range notes {
		model, err := n.getUpdateOutdatedModel(note, authorID)
		if err != nil {
			return fmt.Errorf("note repository: failed to get update outdated model: %v", err)
		}
		notesModel = append(notesModel, model)
	}

	_, err := n.collection.BulkWrite(ctx, notesModel)
	if err != nil {
		return fmt.Errorf("note repository: failed to bulk update notes: %v", err)
	}

	return nil

}

func (n *NoteRepository) getUpdateOutdatedModel(note models.Note, authorID string) (mongo.WriteModel, error) {
	savedNote, err := n.GetNote(note.ExternalID, authorID)
	if err != nil {
		return nil, fmt.Errorf("note repository: failed to get note: %v", err)
	}
	noteNotExist := savedNote == nil
	updatedNote := n.getUpdateNote(note)
	updatedNote["authorId"] = authorID

	if noteNotExist {
		updatedNote["createdAt"] = note.CreatedAt
		return mongo.NewInsertOneModel().SetDocument(updatedNote), nil
	}

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{
			"authorId":   authorID,
			"externalId": note.ExternalID,
			"updatedAt":  bson.M{"$lt": note.UpdatedAt},
		}).
		SetUpdate(bson.M{
			"$set":   updatedNote,
			"$unset": bson.M{"deletedAt": nil},
		}), nil
}

func (n *NoteRepository) DeleteOutdatedNotes(noteIDs []string, authorID string, deletedTime time.Time) error {
	if len(noteIDs) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notesModel := []mongo.WriteModel{}

	for _, noteID := range noteIDs {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{
				"_id":       noteID,
				"authorId":  authorID,
				"updatedAt": bson.M{"$lt": deletedTime},
			}).
			SetUpdate(bson.M{
				"$set":   bson.M{"deletedAt": deletedTime},
				"$unset": bson.M{"updatedAt": deletedTime},
			})
		notesModel = append(notesModel, model)
	}

	_, err := n.collection.BulkWrite(ctx, notesModel)

	if err != nil {
		return fmt.Errorf("note repository: failed to bulk update notes: %v", err)
	}

	return nil

}

type AvailableSpaceInfo struct {
	UsedSpace int64    `bson:"usedSpace"`
	Files     []string `bson:"files"`
}

func (n *NoteRepository) GetUsedSpaceInfo(userID string) (*AvailableSpaceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matchStage := bson.D{{"$match", bson.M{"authorId": userID}}}
	projectStage := bson.D{{
		"$project", bson.M{
			"size":   bson.M{"$bsonSize": "$$ROOT"},
			"images": "$meta.images",
		},
	}}
	groupedByUsedSpaceStage := bson.D{{
		"$group", bson.M{
			"_id":       nil,
			"usedSpace": bson.M{"$sum": "$size"},
			"images":    bson.M{"$addToSet": "$images"},
		},
	}}

	fileteredFilesProjectStage := bson.D{{
		"$project", bson.M{
			"usedSpace": "$usedSpace",
			"images": bson.M{
				"$filter": bson.M{
					"input": "$images",
					"as":    "img",
					"cond":  bson.M{"$ne": bson.A{"$$img", nil}},
				},
			},
		},
	}}

	unwindStage := bson.D{{"$unwind", bson.M{"path": "$images", "preserveNullAndEmptyArrays": false}}}

	groupStage := bson.D{{
		"$group", bson.M{
			"_id":       nil,
			"usedSpace": bson.M{"$last": "$usedSpace"},
			"files":     bson.M{"$push": "$images"},
		},
	}}

	cur, err := n.collection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		projectStage,
		groupedByUsedSpaceStage,
		fileteredFilesProjectStage,
		unwindStage,
		unwindStage,
		groupStage,
	})
	if err != nil {
		return nil, fmt.Errorf("note repository: get used space info: failed to aggregate: %v", err)
	}

	var res AvailableSpaceInfo
	if cur.Next(ctx) {
		err := cur.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("note repository: get used space info: failed to decode: %v", err)
		}
	}

	return &res, nil
}
