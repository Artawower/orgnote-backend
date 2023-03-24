package services

import (
	"fmt"
	"moonbrain/app/models"
	"moonbrain/app/repositories"

	"github.com/rs/zerolog/log"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository}
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

func (u *UserService) FindUser(token string) (*models.PublicUser, error) {
	user, err := u.userRepository.FindUserByToken(token)
	if err != nil {
		return nil, fmt.Errorf("user service: find user: %v", err)
	}
	return mapToPublicUserInfo(user), nil
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
