package services

import (
	"context"
	"errors"
	"strings"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/utils"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, payload *dto.RegisterRequest, config *config.Config) (*dto.TokenResponse, error)
	Login(ctx context.Context, payload *dto.LoginRequest, config *config.Config) (*dto.TokenResponse, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *authService {
	return &authService{repo: repo}
}

func (s *authService) Register(ctx context.Context, payload *dto.RegisterRequest, config *config.Config) (*dto.TokenResponse, error) {
	// check if the user exists
	exists, err := s.repo.IsUserExistByEmail(payload.Email)
	if err != nil {
		return nil, ErrInternal
	}
	if exists {
		return nil, ErrUserExists
	}

	// hash password
	payload.Password, err = auth.HashPassword(payload.Password, config)
	if err != nil {
		return nil, ErrInternal
	}

	// create user with hashed password
	userID, err := s.repo.CreateUser(payload)
	if err != nil {
		return nil, ErrInternal
	}

	return createTokenResponse(config.JWTSecret, userID, config.JWTAccessTokenExpirationInSeconds, config.JWTRefreshTokenExpirationInSeconds)
}

func (s *authService) Login(ctx context.Context, payload *dto.LoginRequest, config *config.Config) (*dto.TokenResponse, error) {
	// проверка существования пользователя
	var user models.User
	var err error
	if strings.Contains(payload.Login, "@") {
		user, err = s.repo.GetUserByEmail(payload.Login)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrUserNotFound
			}
			return nil, ErrInternal
		}
	} else {
		user, err = s.repo.GetUserByNickname(payload.Login)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrUserNotFound
			}
			return nil, ErrInternal
		}
	}

	// проверка пароля
	equal, err := auth.ComparePasswords(payload.Password, user.PasswordHash)
	if err != nil {
		return nil, ErrInternal
	}
	if !equal {
		return nil, ErrUserNotFound
	}

	return createTokenResponse(config.JWTSecret, user.ID, config.JWTAccessTokenExpirationInSeconds, config.JWTRefreshTokenExpirationInSeconds)
}

func createTokenResponse(secret string, userID uint, accessTokenExpirationTime, refreshTokenExpirationTime int) (*dto.TokenResponse, error) {
	accessToken, err := utils.CreateToken(secret, userID, accessTokenExpirationTime)
	if err != nil {
		return nil, ErrInternal
	}

	refreshToken, err := utils.CreateToken(secret, userID, refreshTokenExpirationTime)
	if err != nil {
		return nil, ErrInternal
	}

	return &dto.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
