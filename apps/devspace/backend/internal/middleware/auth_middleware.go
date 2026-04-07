package middleware

import (
	"net/http"
	"strings"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// я бываю невнимателен, проще вынести ключ в константу, чем вспоминать как оно правильно печатается
const UserIdKey = "userID"

func AuthMiddleware(JWTSecret string, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, statusCode := GetUserID(JWTSecret, userRepo, c)
		if statusCode != 0 {
			c.Status(statusCode)
			c.Abort()
			return
		}

		// закидываем userID в контекст
		c.Set(UserIdKey, userID)
		c.Next()
	}
}

// получает из контекста jwt токен и получает пользователя. Возвращает uuid и статус код
func GetUserID(JWTSecret string, userRepo repository.UserRepository, c *gin.Context) (uuid.UUID, int) {
	token := getAccessToken(c)
	if token == "" {
		return uuid.Nil, http.StatusUnauthorized
	}

	userID, err := auth.ValidateJWT([]byte(JWTSecret), token)
	if err != nil {
		return uuid.Nil, http.StatusUnauthorized
	}

	exists, err := userRepo.IsUserExistByID(userID)
	if err != nil {
		return uuid.Nil, http.StatusInternalServerError
	}
	if !exists {
		return uuid.Nil, http.StatusUnauthorized
	}

	return userID, 0
}

func getAccessToken(c *gin.Context) string {
	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if authHeader == "" {
		return ""
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return parts[1]
}
