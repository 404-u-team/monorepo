package repository

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(payload *dto.RegisterRequest) (uint, error)
	IsUserExistByEmail(email string) (bool, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByNickname(login string) (models.User, error)
}

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &userRepository{conn: conn}
}

func (r *userRepository) CreateUser(payload *dto.RegisterRequest) (uint, error) {
	user := models.User{
		Email:        payload.Email,
		Nickname:     payload.Nickname,
		PasswordHash: payload.Password,
	}
	result := r.conn.Create(&user)
	if result.Error != nil {
		log.Println("Ошибка при создании пользователя: ", result.Error)
		return 0, result.Error
	}

	return user.ID, nil
}

func (r *userRepository) IsUserExistByEmail(email string) (bool, error) {
	var exists bool
	err := r.conn.Model(&models.User{}).
		Select("COUNT(*) = 1").
		Where("email = ?", email).
		Find(&exists).Error
	if err != nil {
		log.Println("Ошибка при проверке наличия пользователя: ", err)
		return false, err
	}

	return exists, err
}

func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	result := r.conn.First(&user, "email = ?", email)
	if result.Error != nil {
		log.Println("Ошибка при получении пользователя по email: ", result.Error)
		return user, result.Error
	}

	return user, nil
}

func (r *userRepository) GetUserByNickname(nickname string) (models.User, error) {
	var user models.User

	result := r.conn.First(&user, "nickname = ?", nickname)
	if result.Error != nil {
		log.Println("Ошибка при получении пользователя по nickname: ", result.Error)
		return user, result.Error
	}

	return user, nil
}
