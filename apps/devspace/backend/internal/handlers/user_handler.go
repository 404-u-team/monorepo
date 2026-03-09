package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) GetMe(c *gin.Context) {
	// получение ID пользователя из контекста
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetMe(userID)
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

func (h *userHandler) UpdateMe(c *gin.Context) {
	// получение ID пользователя из контекста
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// получение payload
	var payload dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.userService.UpdateMe(userID, &payload)
	if err != nil {
		if errors.Is(err, services.ErrEmptyPayload) {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, services.ErrUserNotFound) {
			c.Status(http.StatusUnauthorized)
			return
		}
		if errors.Is(err, services.ErrUserConflict) {
			c.JSON(http.StatusConflict, err.Error())
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
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
