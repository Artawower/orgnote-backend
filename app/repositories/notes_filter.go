package repositories

import (
	"moonbrain/app/models"

	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
)

func addDeletedFilter(filter bson.M, modelFilter models.NoteFilter) {
	if modelFilter.IncludeDeleted == nil {
		modelFilter.IncludeDeleted = new(bool)
	}

	includeDeleted := *modelFilter.IncludeDeleted
	orQuery := []bson.M{{"deletedAt": bson.M{"$exists": includeDeleted}}}

	if !includeDeleted {
		orQuery = append(orQuery, bson.M{"deletedAt": bson.M{"$eq": nil}})
	}

	filter["$or"] = orQuery
}

func addPublishedFilter(filter bson.M, modelFilter models.NoteFilter) {
	if modelFilter.Published == nil {
		return
	}
	orQuery := []bson.M{{"meta.published": bson.M{"$eq": *modelFilter.Published}}}

	if !*modelFilter.Published {
		orQuery = append(orQuery, bson.M{"meta.published": bson.M{"$exists": false}})
	}
	filter["$or"] = orQuery
}

func addAuthorIdFilter(filter bson.M, modelFilter models.NoteFilter) {
	if modelFilter.UserID == nil {
		return
	}

	filter["authorId"] = *modelFilter.UserID
}

func addUpdatedTimeFilter(filter bson.M, modelFilter models.NoteFilter) {
	if modelFilter.From == nil {
		return
	}
	filter["updatedAt"] = bson.M{"$gte": *modelFilter.From}
}

func addSearchFilter(filter *bson.M, modelFilter models.NoteFilter) {
	if modelFilter.SearchText != nil && *modelFilter.SearchText != "" {
		f := *filter
		f["$text"] = bson.D{bson.E{Key: "$search", Value: *modelFilter.SearchText}}
	}
}

var filterBuilders = []func(filter bson.M, modelFilter models.NoteFilter){
	addDeletedFilter,
	addPublishedFilter,
	addAuthorIdFilter,
	addUpdatedTimeFilter,
}

func getNotesFilter(modelFilter models.NoteFilter) bson.M {
	filter := bson.M{}

	funk.ForEach(filterBuilders, func(builder func(f bson.M, modelFilter models.NoteFilter)) {
		builder(filter, modelFilter)
	})
	return filter
}
