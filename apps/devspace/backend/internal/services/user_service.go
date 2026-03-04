package services

import "github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"

type UserService interface {
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *userService {
	return &userService{repo: repo}
}
