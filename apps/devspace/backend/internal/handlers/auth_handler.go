package handlers

import (
	"net/http"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService services.AuthService
	config      *config.Config
}

func NewAuthHandler(userService services.AuthService, config *config.Config) *authHandler {
	return &authHandler{
		authService: userService,
		config:      config,
	}
}

func (h *authHandler) Register(c *gin.Context) {
	var payload dto.RegisterRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenResponse, err := h.authService.Register(c, &payload, h.config)
	if err != nil {
		// TODO check for different error response -> differenct error code
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setTokenIntoCookie(c, tokenResponse.RefreshToken, h.config.JWTRefreshTokenExpirationInSeconds)
	c.JSON(http.StatusCreated, gin.H{"access_token": tokenResponse.AccessToken})
}

// TODO: Добавить /refresh

func setTokenIntoCookie(c *gin.Context, token string, expirationTime int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token",
		token,
		expirationTime, // время жизни внутри куки
		"/",
		"",
		false, // когда будем использовать https поставить на true
		true,
	)
}
