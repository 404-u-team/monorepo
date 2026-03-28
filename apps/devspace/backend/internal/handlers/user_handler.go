package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/middleware"
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

	meResponse, err := h.userService.UpdateMe(userID, &payload)
	if err != nil {
		if errors.Is(err, services.ErrEmptyPayload) || errors.Is(err, services.ErrSkillNotFound) || errors.Is(err, services.ErrSkillIsNotRoot) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrUserNotFound) {
			c.Status(http.StatusUnauthorized)
			return
		}
		if errors.Is(err, services.ErrUserConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, meResponse)
}

func (h *userHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("userID")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	getMeResponse, err := h.userService.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, getMeResponse)
}

func getUserId(c *gin.Context) (uuid.UUID, error) {
	userIDAny, ok := c.Get(middleware.UserIdKey)
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

func (h *userHandler) GetUsersByParams(c *gin.Context) {
	var req dto.GetUsersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("Binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid parameters",
			"details": err.Error(),
		})
		return
	}

	// Парсим main_role (если передан)
	var mainRoleUUID *uuid.UUID
	if req.MainRole != nil && *req.MainRole != "" {
		id, err := uuid.Parse(*req.MainRole)
		if err != nil {
			log.Printf("Invalid main_role format: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid main_role format",
				"details": "UUID expected",
			})
			return
		}
		mainRoleUUID = &id
	}

	log.Printf("Request: start_at=%v, limit=%v, username=%v, main_role=%v, skills=%v",
		req.StartAt, req.Limit, req.Username, mainRoleUUID, req.Skills)

	// Получаем публичные профили с деревом навыков
	profiles, err := h.userService.GetUsersPublicProfiles(
		req.StartAt,
		req.Limit,
		req.Username,
		mainRoleUUID,
		req.Skills,
	)

	if err != nil {
		log.Printf("Error getting users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка получения пользователей",
		})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

// TODO: Прописать журналирование для остальных 500 в этом хэндлере. Мы же потом будем с недоумением смотреть на код.
