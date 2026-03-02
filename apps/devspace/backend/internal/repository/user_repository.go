package repository

import (
	"context"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, payload *dto.RegisterRequest) (int, error)
	IsUserExistByEmail(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &userRepository{conn: conn}
}

func (u *userRepository) CreateUser(ctx context.Context, payload *dto.RegisterRequest) (int, error) {
	return 0, nil
}

func (u *userRepository) IsUserExistByEmail(ctx context.Context, email string) (bool, error) {
	return true, nil
}
