package handlers

import (
	"errors"
	"log"
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

func NewAuthHandler(authService services.AuthService, config *config.Config) *authHandler {
	return &authHandler{
		authService: authService,
		config:      config,
	}
}

func (h *authHandler) Register(c *gin.Context) {
	var payload dto.RegisterRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	tokenResponse, err := h.authService.Register(&payload, h.config)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			c.Status(http.StatusConflict)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	setTokenIntoCookie(c, tokenResponse.RefreshToken, h.config.JWTRefreshTokenExpirationInSeconds, h.config.AllowAnyOrigin)
	c.JSON(http.StatusCreated, gin.H{"access_token": tokenResponse.AccessToken})
}

func (h *authHandler) Login(c *gin.Context) {
	var payload dto.LoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	tokenResponse, err := h.authService.Login(&payload, h.config)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	setTokenIntoCookie(c, tokenResponse.RefreshToken, h.config.JWTRefreshTokenExpirationInSeconds, h.config.AllowAnyOrigin)
	c.JSON(http.StatusOK, gin.H{"access_token": tokenResponse.AccessToken})
}

func (h *authHandler) Refresh(c *gin.Context) {
	tokenResponse, err := h.authService.Refresh(c, h.config)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	setTokenIntoCookie(c, tokenResponse.RefreshToken, h.config.JWTRefreshTokenExpirationInSeconds, h.config.AllowAnyOrigin)
	c.JSON(http.StatusOK, gin.H{"access_token": tokenResponse.AccessToken})
}

func setTokenIntoCookie(c *gin.Context, token string, expirationTime int, allowAnyOrigin bool) {
	sameSite := http.SameSiteLaxMode
	if allowAnyOrigin {
		sameSite = http.SameSiteNoneMode
	}
	c.SetSameSite(sameSite)
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
