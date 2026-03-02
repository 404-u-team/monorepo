package services

import (
	"context"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	authpb "github.com/404-u-team/monorepo/apps/devspace/backend/services/proto/auth/v1"
)

type AuthService interface {
	Register(ctx context.Context, payload *authpb.RegisterRequest) (int, error)
	// Login(ctx context.Context, payload *authpb.LoginRequest) (int, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *authService {
	return &authService{repo: repo}
}

func (s *authService) Register(ctx context.Context, payload *authpb.RegisterRequest) (int, error) {
	// check if the user exists
	exists, err := s.repo.IsUserExistByEmail(ctx, payload.Email)
	if err != nil {
		return -1, ErrInternal
	}
	if exists {
		return -1, ErrUserExists
	}

	// hash password
	payload.Password, err = auth.HashPassword(payload.Password)
	if err != nil {
		return -1, ErrInternal
	}

	// create user with hashed password
	userId, err := s.repo.CreateUser(ctx, payload)
	if err != nil {
		return -1, ErrInternal
	}

	return userId, nil
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
