package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type projectHandler struct {
	projectService services.ProjectService
	config         *config.Config
}

func NewProjectHandler(projectService services.ProjectService, config *config.Config) *projectHandler {
	return &projectHandler{
		projectService: projectService,
		config:         config,
	}
}

func (h *projectHandler) CreateProject(c *gin.Context) {
	var payload dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	userIDAny, ok := c.Get("userID")
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	userID, ok := userIDAny.(uuid.UUID)
	if !ok {
		log.Println("Ошибка при конвертировании userID в UUID")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := h.projectService.CreateProject(&payload, userID); err != nil {
		if errors.Is(err, services.ErrProjectConflict) {
			c.Status(http.StatusConflict)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
