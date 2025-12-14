package repositories

import (
	"context"
	"errors"
	"fmt"
	"orgnote/app/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db:         db,
		collection: db.Collection("users"),
	}
}

func (u *UserRepository) CreateOrGet(user models.User) (*models.User, error) {
	foundUser, err := u.GetUser(&user)

	if foundUser != nil {
		updatedUser, err := u.UpdateAuthInfo(user)
		if err != nil {
			return nil, fmt.Errorf("user repository: create or update user: update auth info: %v", err)
		}
		return updatedUser, nil
	}
	createdUser, err := u.Create(user)
	if err != nil {
		return nil, fmt.Errorf("user repository: create or update user: create user: %v", err)
	}

	return createdUser, nil
}

func (u *UserRepository) UpdateAuthInfo(user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"externalId": user.ExternalID, "provider": user.Provider}

	_, err := u.collection.UpdateOne(ctx, filter, bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "token", Value: user.Token},
			bson.E{Key: "refreshToken", Value: user.RefreshToken},
			bson.E{Key: "tokenExpiration", Value: user.TokenExpirationDate},
			bson.E{Key: "profileUrl", Value: user.ProfileURL},
		}},
	})

	if err != nil {
		return nil, fmt.Errorf("user repository: update user: update one user: %v", err)
	}

	updatedUser, err := u.GetUser(&user)

	if err != nil {
		return nil, fmt.Errorf("user repository: update user: get user: %v", err)
	}

	return updatedUser, nil
}

func (u *UserRepository) Create(user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user.ID = primitive.NewObjectID()
	_, err := u.collection.InsertOne(ctx, user)

	if err != nil {
		return nil, fmt.Errorf("user repository: create user: insert one user: %v", err)
	}

	createdUser, err := u.GetUser(&user)

	if err != nil {
		return nil, fmt.Errorf("user repository: create user: get user: %v", err)
	}
	return createdUser, nil
}

func (u *UserRepository) GetUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"externalId": user.ExternalID, "provider": user.Provider}
	err := u.collection.FindOne(ctx, filter).Decode(user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("user repository: get user: find one user: %v", err)
	}
	return user, nil
}

func (u *UserRepository) GetByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("user repository: get by id: convert id: %v", err)
	}
	filter := bson.M{"_id": objID}
	user := models.User{}
	err = u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user repository: find user by id: find one user: %v", err)
	}
	return &user, nil
}

func (u *UserRepository) GetUsersByIDs(userIDs []string) ([]models.User, error) {
	objectUserIDs := make([]primitive.ObjectID, len(userIDs))
	for i, id := range userIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("user repository: get users by ids: convert id - %s: %v", id, err)
		}
		objectUserIDs[i] = objID
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": bson.M{"$in": objectUserIDs}}
	users := []models.User{}
	cur, err := u.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("user repository: get users by ids: find users: %v", err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("user repository: get users by ids: decode user: %v", err)
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("user repository: get users by ids: cursor error: %v", err)
	}
	return users, nil
}

func (u *UserRepository) FindUserByToken(token string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"$or": bson.A{
		bson.M{"token": token},
		bson.M{"apiTokens": bson.M{"$elemMatch": bson.M{"token": token}}},
	}}
	user := models.User{}
	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user repository: find user by token: find one user: %v", err)
	}
	return &user, nil
}

func (u *UserRepository) GetAPITokens(userID string) ([]models.APIToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("user repository: get api tokens: convert user id: %v", err)
	}

	filter := bson.M{"_id": userObjID}
	user := models.User{}
	err = u.collection.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []models.APIToken{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("user repository: get api tokens: find one user: %v", err)
	}
	return user.APITokens, nil
}

func (u *UserRepository) CreateAPIToken(user *models.User) (*models.APIToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": user.ID}
	token := uuid.New()
	accessToken := models.APIToken{
		ID:          primitive.NewObjectID(),
		Permissions: "w",
		Token:       token.String(),
	}
	update := bson.M{"$push": bson.M{"apiTokens": accessToken}}

	_, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("user repository: create api token: update one user: %v", err)
	}
	return &accessToken, nil
}

// Delete user API token from list of tokens
func (u *UserRepository) DeleteAPIToken(user *models.User, tokenID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": user.ID}
	id, err := primitive.ObjectIDFromHex(tokenID)
	if err != nil {
		return fmt.Errorf("user repository: delete api token: convert token id: %v", err)
	}
	update := bson.M{"$pull": bson.M{"apiTokens": bson.M{"_id": id}}}

	_, err = u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("user repository: delete api token: update one user: %v", err)
	}

	return nil
}

func (u *UserRepository) GetAll() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	users := []models.User{}
	cur, err := u.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("user repository: get all: find users: %v", err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("user repository: get all: decode user: %v", err)
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("user repository: get all: cursor error: %v", err)
	}
	return users, nil
}

func (u *UserRepository) UpdateSpaceLimitInfo(usedID string, usedSpace *int64, spaceLimit *int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(usedID)

	if err != nil {
		return fmt.Errorf("note repository: update used space: convert id: %v", err)
	}

	filter := bson.M{"_id": objID}

	updatedModel := bson.M{}
	if usedSpace != nil {
		updatedModel["usedSpace"] = usedSpace
	}
	if spaceLimit != nil {
		updatedModel["spaceLimit"] = spaceLimit
	}
	update := bson.M{"$set": updatedModel}

	_, err = u.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("note repository: update used space: failed to update: %v", err)
	}

	return nil
}

func (u *UserRepository) SetActivationKey(userID string, activationKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return fmt.Errorf("note repository: set active status: convert id: %v", err)
	}

	filter := bson.M{"_id": objID}

	update := bson.M{"$set": bson.M{"active": activationKey}}

	_, err = u.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("note repository: set active status: failed to update: %v", err)
	}

	return nil
}

func (u *UserRepository) DeleteUser(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return fmt.Errorf("note repository: delete user: convert id: %v", err)
	}

	filter := bson.M{"_id": objID}

	_, err = u.collection.DeleteOne(ctx, filter)

	if err != nil {
		return fmt.Errorf("note repository: delete user: failed to delete: %v", err)
	}

	return nil
}
