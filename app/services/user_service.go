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
	subscriptionAPI *infrastructure.SubscriptionAPI
}

func NewUserService(userRepository *repositories.UserRepository, subscriptionAPI *infrastructure.SubscriptionAPI) *UserService {
	return &UserService{userRepository, subscriptionAPI}
}

func (u *UserService) FindOrCreate(user models.User) (*models.User, error) {
	log.Info().Str("provider", user.Provider).Str("externalId", user.ExternalID).Msg("Find or create user")
	createdUser, err := u.userRepository.CreateOrGet(user)
	if err != nil {
		return nil, fmt.Errorf("user service: find or create: %v", err)
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
	return nil
}

func (u *UserService) Subscribe(user *models.User, token string, emailAddress *string) error {
	var email *types.Email
	if emailAddress != nil {
		email = (*types.Email)(emailAddress)
	}
	var externalEmail *types.Email
	if user.Email != "" {
		email = (*types.Email)(&user.Email)
	}
	data, err := u.subscriptionAPI.ActivateSubscription(subscription.SubscriptionActivation{
		Key:              token,
		Email:            email,
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
