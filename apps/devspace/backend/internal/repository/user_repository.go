package repository

import (
	"errors"
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(payload *dto.RegisterRequest) (uuid.UUID, error)
	IsUserExistByEmail(email string) (bool, error)
	IsUserExistByID(id uuid.UUID) (bool, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByNickname(login string) (models.User, error)
	CheckUserIsAdmin(id string) (bool, error)
}

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &userRepository{conn: conn}
}

func (r *userRepository) CreateUser(payload *dto.RegisterRequest) (uuid.UUID, error) {
	user := models.User{
		Email:        payload.Email,
		Nickname:     payload.Nickname,
		PasswordHash: payload.Password,
	}
	result := r.conn.Create(&user)
	if result.Error != nil {
		log.Println("Ошибка при создании пользователя: ", result.Error)
		return uuid.Nil, result.Error
	}

	return user.ID, nil
}

func (r *userRepository) IsUserExistByEmail(email string) (bool, error) {
	var count int64
	err := r.conn.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		log.Println("Ошибка при проверке наличия пользователя: ", err)
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) IsUserExistByID(id uuid.UUID) (bool, error) {
	var count int64
	err := r.conn.Model(&models.User{}).
		Where("id = ?", id).
		Count(&count).Error
	if err != nil {
		log.Println("Ошибка при проверке наличия пользователя: ", err)
		return false, err
	}
	return count > 0, nil
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

func (r *userRepository) CheckUserIsAdmin(id string) (bool, error) {
	var user models.User

	result := r.conn.Select("is_admin").Where("id = ?", id).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("Пользователь не найден по id ", id)
			return false, nil
		}
		log.Println("Ошибка при получении прав юзера по id ", result.Error)
		return false, result.Error
	}

	return user.IsAdmin, nil
}
