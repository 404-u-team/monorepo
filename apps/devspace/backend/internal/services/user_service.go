package services

import (
	"errors"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetMe(userID uuid.UUID) (*dto.GetMeResponse, error)
	UpdateMe(userID uuid.UUID, updateRequest *dto.UpdateUserRequest) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) GetMe(userID uuid.UUID) (*dto.GetMeResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}

	userSkills, err := s.repo.GetUserSkills(userID)
	if err != nil {
		return nil, ErrInternal
	}

	getMeResponse := dto.GetMeResponse{
		ID:        userID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarUri: "",
		Bio:       user.Bio,
		Skills:    userSkills,
	}
	return &getMeResponse, nil
}

func (s *userService) UpdateMe(userID uuid.UUID, updateRequest *dto.UpdateUserRequest) error {
	if updateRequest.Nickname == nil && updateRequest.Bio == nil {
		return ErrEmptyPayload
	}

	exists, err := s.repo.IsUserExistByID(userID)
	if err != nil {
		return ErrInternal
	}
	if !exists {
		return ErrUserNotFound
	}

	err = s.repo.UpdateUserByID(userID, updateRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrUserConflict
		}

		return ErrInternal
	}

	return nil
}
