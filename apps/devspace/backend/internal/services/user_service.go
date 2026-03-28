package services

import (
	"errors"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetMe(userID uuid.UUID) (*dto.PrivateUserProfile, error)
	UpdateMe(userID uuid.UUID, updateRequest *dto.UpdateUserRequest) (*dto.PrivateUserProfile, error)
	GetUserByID(userID uuid.UUID) (*dto.PublicUserProfile, error)
	GetUsersByParams(startAt *uint, limit *uint, username *string, mainRole *uuid.UUID, skills *dto.UUIDSlice) ([]models.User, error)
	GetUsersPublicProfiles(
		startAt, limit *uint,
		username *string,
		mainRole *uuid.UUID,
		skills *dto.UUIDSlice,
	) ([]dto.PublicUserProfile, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) GetMe(userID uuid.UUID) (*dto.PrivateUserProfile, error) {
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

	privateUserProfile := dto.PrivateUserProfile{
		ID:        userID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarUri: "",
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
		Skills:    userSkills,
	}
	return &privateUserProfile, nil
}

func (s *userService) UpdateMe(userID uuid.UUID, updateRequest *dto.UpdateUserRequest) (*dto.PrivateUserProfile, error) {
	if updateRequest.Nickname == nil && updateRequest.Bio == nil && updateRequest.AvatarUrl == nil{
		return nil, ErrEmptyPayload
	}

	exists, err := s.repo.IsUserExistByID(userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	err = s.repo.UpdateUserByID(userID, updateRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrUserConflict
		}

		return nil, ErrInternal
	}

	getMeResponse, err := s.GetMe(userID)
	if err != nil {
		return nil, err
	}

	return getMeResponse, nil
}

func (s *userService) GetUserByID(userID uuid.UUID) (*dto.PublicUserProfile, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}

	userSkills, err := s.repo.GetUserSkills(user.ID)
	if err != nil {
		return nil, ErrInternal
	}

	getMeResponse := dto.PublicUserProfile{
		ID:        user.ID,
		Nickname:  user.Nickname,
		MainRole:  user.MainRole,
		AvatarUri: user.AvatarUrl,
		Bio:       user.Bio,
		Skills:    userSkills,
	}

	return &getMeResponse, nil
}

func (s *userService) GetUsersByParams(
	startAt, limit *uint,
	username *string,
	mainRole *uuid.UUID,
	skills *dto.UUIDSlice,
) ([]models.User, error) {
	return s.repo.GetUsersByParams(startAt, limit, username, mainRole, skills)
}

func (s *userService) GetUsersPublicProfiles(
	startAt, limit *uint,
	username *string,
	mainRole *uuid.UUID,
	skills *dto.UUIDSlice,
) ([]dto.PublicUserProfile, error) {

	// Получаем пользователей с навыками (Preload уже подгрузил Skills)
	users, err := s.repo.GetUsersByParams(startAt, limit, username, mainRole, skills)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []dto.PublicUserProfile{}, nil
	}

	// Конвертируем в публичные профили с деревом навыков
	profiles := make([]dto.PublicUserProfile, len(users))
	for i, user := range users {
		// Строим дерево навыков для пользователя
		skillTree, dbError := s.repo.GetUserSkills(user.ID)

		if dbError != nil {
			return nil, dbError
		}

		profiles[i] = dto.PublicUserProfile{
			ID:        user.ID,
			Nickname:  user.Nickname,
			MainRole:  user.MainRole,
			AvatarUri: user.AvatarUrl,
			Bio:       user.Bio,
			Skills:    skillTree,
		}
	}

	return profiles, nil
}
