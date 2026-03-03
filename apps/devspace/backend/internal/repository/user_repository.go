package repository

import (
	"context"
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, payload *dto.RegisterRequest) (uint, error)
	IsUserExistByEmail(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &userRepository{conn: conn}
}

func (r *userRepository) CreateUser(ctx context.Context, payload *dto.RegisterRequest) (uint, error) {
	user := models.User{
		Email:        payload.Email,
		Nickname:     payload.Nickname,
		PasswordHash: payload.Password,
	}
	result := r.conn.Create(&user)
	if result.Error != nil {
		log.Printf("Ошибка при создании пользователя: %v\n", result.Error)
		return 0, result.Error
	}

	return user.ID, nil
}

func (r *userRepository) IsUserExistByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.conn.Model(&models.User{}).
		Select("COUNT(*) = 1").
		Where("email = ?", email).
		Find(&exists).Error
	if err != nil {
		log.Printf("Ошибка при проверке наличия пользователя: %v\n", err)
	}

	return exists, err
}
