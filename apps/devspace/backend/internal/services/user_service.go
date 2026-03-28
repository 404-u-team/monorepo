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
	GetUsersByParams(startAt *uint, limit *uint, username *string, mainRole *uuid.UUID, skills *dto.UUIDSlice) ([]dto.PublicUserProfile, error)
}

type userService struct {
	userRepo  repository.UserRepository
	skillRepo repository.SkillRepository
}

func NewUserService(userRepo repository.UserRepository, skillRepo repository.SkillRepository) *userService {
	return &userService{userRepo: userRepo, skillRepo: skillRepo}
}

func (s *userService) GetMe(userID uuid.UUID) (*dto.PrivateUserProfile, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}

	var mainRole *models.SkillCategory
	if user.MainRole != nil {
		mainRole, err = s.skillRepo.GetSkillByID(*user.MainRole)
		if err != nil {
			return nil, ErrInternal
		}
	}

	userSkills, err := s.userRepo.GetUserSkills(userID)
	if err != nil {
		return nil, ErrInternal
	}

	privateUserProfile := dto.PrivateUserProfile{
		ID:        userID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		MainRole:  nil,
		AvatarUrl: "",
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
		Skills:    userSkills,
	}
	if user.MainRole != nil {
		privateUserProfile.MainRole = mainRole
	}
	return &privateUserProfile, nil
}

func (s *userService) UpdateMe(userID uuid.UUID, updateRequest *dto.UpdateUserRequest) (*dto.PrivateUserProfile, error) {
	if updateRequest.Nickname == nil && updateRequest.Bio == nil && updateRequest.AvatarUrl == nil && !updateRequest.MainRole.IsSet {
		return nil, ErrEmptyPayload
	}

	exists, err := s.userRepo.IsUserExistByID(userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	err = s.userRepo.UpdateUserByID(userID, updateRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMainRoleNotFound
		}
		if errors.Is(err, gorm.ErrInvalidValue) {
			return nil, ErrMainRoleIsNotRoot
		}
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
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}

	var mainRole *models.SkillCategory
	if user.MainRole != nil {
		mainRole, err = s.skillRepo.GetSkillByID(*user.MainRole)
		if err != nil {
			return nil, ErrInternal
		}
	}

	userSkills, err := s.userRepo.GetUserSkills(user.ID)
	if err != nil {
		return nil, ErrInternal
	}

	publicUserProfile := dto.PublicUserProfile{
		ID:        user.ID,
		Nickname:  user.Nickname,
		MainRole:  nil,
		AvatarUrl: user.AvatarUrl,
		Bio:       user.Bio,
		Skills:    userSkills,
	}
	if user.MainRole != nil {
		publicUserProfile.MainRole = mainRole
	}

	return &publicUserProfile, nil
}

func (s *userService) GetUsersByParams(
	startAt, limit *uint,
	username *string,
	mainRole *uuid.UUID,
	skills *dto.UUIDSlice,
) ([]dto.PublicUserProfile, error) {
	return s.userRepo.GetUsersByParams(startAt, limit, username, mainRole, skills)
}
