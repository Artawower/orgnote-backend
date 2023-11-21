package services

import (
	"fmt"
	"orgnote/app/infrastructure"
	subscription "orgnote/app/infrastructure/generated"
	"orgnote/app/models"
	"orgnote/app/repositories"

	"github.com/davecgh/go-spew/spew"
	"github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"
)

type UserService struct {
	userRepository  *repositories.UserRepository
	noteRepository  *repositories.NoteRepository
	subscriptionAPI *infrastructure.SubscriptionAPI
}

func NewUserService(userRepository *repositories.UserRepository, noteRepository *repositories.NoteRepository, subscriptionAPI *infrastructure.SubscriptionAPI) *UserService {
	return &UserService{userRepository, noteRepository, subscriptionAPI}
}

func (u *UserService) Login(user models.User) (*models.User, error) {
	log.Info().Msgf("Login user: %v", user)
	createdUser, err := u.userRepository.CreateOrGet(user)
	if err != nil {
		return nil, fmt.Errorf("user service: login: %v", err)
	}
	return createdUser, nil
}

func (u *UserService) GetAPITokens(userID string) ([]models.APIToken, error) {
	tokens, err := u.userRepository.GetAPITokens(userID)
	if err != nil {
		return nil, fmt.Errorf("user service: get: %v", err)
	}
	return tokens, nil
}

func (u *UserService) FindUser(token string) (*models.UserPersonalInfo, error) {
	user, err := u.userRepository.FindUserByToken(token)
	if err != nil {
		return nil, fmt.Errorf("user service: find user: %v", err)
	}
	return mapToUserPersonalInfo(user), nil
}

func (u *UserService) CreateToken(user *models.User) (*models.APIToken, error) {
	token, err := u.userRepository.CreateAPIToken(user)
	if err != nil {
		return nil, fmt.Errorf("user service: create token: %v", err)
	}
	return token, nil
}

func (u *UserService) DeleteToken(user *models.User, tokenID string) error {
	err := u.userRepository.DeleteAPIToken(user, tokenID)
	if err != nil {
		return fmt.Errorf("user service: delete token: %v", err)
	}
	return nil
}

func (u *UserService) DeleteUser(user *models.User) error {
	err := u.userRepository.DeleteUser(user.ID.Hex())
	if err != nil {
		return fmt.Errorf("user service: delete user: %v", err)
	}

	err = u.noteRepository.DeleteUserNotes(user.ID.Hex())
	if err != nil {
		return fmt.Errorf("user service: delete user notes: %v", err)
	}
	return nil
}

func (u *UserService) Subscribe(user *models.User, token string, email *string) error {
	// TODO: master transaction with context
	var externalEmail *types.Email
	if user.Email != "" {
		externalEmail = (*types.Email)(&user.Email)
	}
	data, err := u.subscriptionAPI.ActivateSubscription(subscription.SubscriptionActivation{
		Key:              token,
		Email:            (*types.Email)(email),
		ExternalId:       user.ExternalID,
		ExternalEmail:    externalEmail,
		ExternalProvider: &user.Provider,
	})
	spew.Dump(data)
	if err != nil {
		return fmt.Errorf("user service: subscribe: activate subscription %v", err)
	}

	err = u.userRepository.SetActivationKey(user.ID.Hex(), token)
	if err != nil {
		return fmt.Errorf("user service: subscribe: set active status: %v", err)
	}

	spaceLimit := int64(*data.SpaceLimit)

	err = u.userRepository.UpdateSpaceLimitInfo(user.ID.Hex(), nil, &spaceLimit)
	if err != nil {
		return fmt.Errorf("user service: subscribe: update space limit info: %v", err)
	}
	return nil
}

func (u *UserService) calculateUserUsedSpace(user *models.User) error {
	return nil
}
