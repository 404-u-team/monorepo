package grpc

import (
	"context"
	"errors"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/service"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/utils"
	authpb "github.com/404-u-team/monorepo/apps/devspace/backend/services/proto/auth/v1"
)

type authServer struct {
	authpb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthServer(authService service.AuthService) *authServer {
	return &authServer{authService: authService}
}

func (s *authServer) Register(context context.Context, req *authpb.RegisterRequest) (*authpb.TokenResponse, error) {
	userId, err := s.authService.Register(context, req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return nil, service.ErrUserExists
		}
		if errors.Is(err, service.ErrInternal) {
			return nil, service.ErrInternal
		}
	}

	accessToken, err := utils.CreateToken(config.Envs.JWTSecret, userId)
	if err != nil {
		return nil, service.ErrInternal
	}

	return &authpb.TokenResponse{AccessToken: accessToken, RefreshToken: "2"}, nil
}

func (s *authServer) Login(context context.Context, req *authpb.LoginRequest) (*authpb.TokenResponse, error) {
	userId, err := s.authService.Login(context, req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return nil, service.ErrUserExists
		}
		if errors.Is(err, service.ErrInternal) {
			return nil, service.ErrInternal
		}
	}

	accessToken, err := utils.CreateToken(config.Envs.JWTSecret, userId)
	if err != nil {
		return nil, service.ErrInternal
	}

	return &authpb.TokenResponse{AccessToken: accessToken, RefreshToken: "2"}, nil
}

func (s *authServer) ValidateToken(context context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	return &authpb.ValidateTokenResponse{UserId: 18, Email: "test@mail.com", Valid: true}, nil
}
