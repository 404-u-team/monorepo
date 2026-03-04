package services

import (
	"errors"
	"log"
	"strings"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(payload *dto.RegisterRequest, config *config.Config) (*dto.TokenResponse, error)
	Login(payload *dto.LoginRequest, config *config.Config) (*dto.TokenResponse, error)
	Refresh(c *gin.Context, config *config.Config) (*dto.TokenResponse, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *authService {
	return &authService{repo: repo}
}

func (s *authService) Register(payload *dto.RegisterRequest, config *config.Config) (*dto.TokenResponse, error) {
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

func (s *authService) Login(payload *dto.LoginRequest, config *config.Config) (*dto.TokenResponse, error) {
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

func (s *authService) Refresh(c *gin.Context, config *config.Config) (*dto.TokenResponse, error) {
	// получение токена из куки
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return nil, ErrUnauthorized
	}

	// проверка токена
	userID, err := auth.ValidateJWT([]byte(config.JWTSecret), refreshToken)
	if err != nil {
		log.Println("Ошибка при валидации jwt токена в /refresh: ", err)
		return nil, ErrUnauthorized
	}

	// создание новых токенов
	return createTokenResponse(config.JWTSecret, userID, config.JWTAccessTokenExpirationInSeconds, config.JWTRefreshTokenExpirationInSeconds)
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
