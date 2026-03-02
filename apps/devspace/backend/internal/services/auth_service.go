package services

import (
	"context"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, payload *dto.RegisterRequest, config *config.Config) (*dto.TokenResponse, error)
	// Login(ctx context.Context, payload *authpb.LoginRequest) (int, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *authService {
	return &authService{repo: repo}
}

func (s *authService) Register(ctx context.Context, payload *dto.RegisterRequest, config *config.Config) (*dto.TokenResponse, error) {
	// check if the user exists
	exists, err := s.repo.IsUserExistByEmail(ctx, payload.Email)
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
	userID, err := s.repo.CreateUser(ctx, payload)
	if err != nil {
		return nil, ErrInternal
	}

	accessToken, err := utils.CreateToken(config.JWTSecret, userID, config)
	if err != nil {
		return nil, ErrInternal
	}

	return &dto.TokenResponse{AccessToken: accessToken, RefreshToken: "2"}, nil
}

// func (s *authService) Login(ctx context.Context, payload *authpb.LoginRequest) (int, error) {
// 	// check if the user exists
// 	u, err := s.repo.GetUserByEmail(ctx, payload.Email)
// 	if err != nil {
// 		return -1, ErrInternal
// 	}

// 	// compare password
// 	if !auth.ComparePasswords(u.Password, payload.Password) {
// 		return -1, ErrUserNotFound
// 	}

// 	return u.ID, nil
// }
