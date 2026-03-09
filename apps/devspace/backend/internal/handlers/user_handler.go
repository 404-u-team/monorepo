package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandler struct {
	userService services.UserService
	config      *config.Config
}

func NewUserHandler(userService services.UserService, config *config.Config) *userHandler {
	return &userHandler{
		userService: userService,
		config:      config,
	}
}

func (h *userHandler) Me(c *gin.Context) {
	// получение ID пользователя из контекста
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	user, err := h.userService.Me(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func getUserId(c *gin.Context) (uuid.UUID, error) {
	userIDAny, ok := c.Get("userID")
	if !ok {
		return uuid.Nil, services.ErrUnauthorized
	}

	userID, ok := userIDAny.(uuid.UUID)
	if !ok {
		log.Println("Ошибка при конвертировании userID в UUID")
		return uuid.Nil, services.ErrUnauthorized
	}

	return userID, nil
}
