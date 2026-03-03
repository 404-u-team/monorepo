package handlers

import (
	"net/http"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, gin.H{"status": "authorized"})
}
