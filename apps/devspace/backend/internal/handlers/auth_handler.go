package handlers

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	userService *services.AuthService
}

func NewAuthHandler(userService *services.AuthService) *authHandler {
	return &authHandler{
		userService: userService,
	}
}

func (h *authHandler) Register(ctx *gin.Context) {
	// получаем данные

	// отправляем их в services.Register

	// получаем результат и отправляем их
}
