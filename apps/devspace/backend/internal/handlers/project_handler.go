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

	project, err := h.projectService.CreateProject(&payload, userID)
	if err != nil {
		if errors.Is(err, services.ErrProjectConflict) {
			c.Status(http.StatusConflict)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *projectHandler) GetProjects(c *gin.Context) {
	var query dto.GetProjectsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Println("Ошибка при парсинге query: ", err)
		c.Status(http.StatusBadRequest)
		return
	}

	projects, err := h.projectService.GetProjects(&query)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *projectHandler) GetProjectByID(c *gin.Context) {
	projectIDStr := c.Param("projectID")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *projectHandler) UpdateProjectByID(c *gin.Context) {
	// получение projectID из параметров
	projectID, err := getProjectID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrProjectNotFound.Error()})
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// является ли пользователь владельцем данного проекта
	if projectID != userID {
		c.Status(http.StatusForbidden)
		return
	}

	// получение payload
	var payload dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	if payload.Title == nil && payload.Description == nil && payload.Status == nil {
		log.Println("Все поля пустые, нечего изменять")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля пустые, нечего изменять"})
		return
	}

	err = h.projectService.UpdateProjectById(projectID, &payload)
	if err != nil {
		if errors.Is(err, services.ErrProjectConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		if errors.Is(err, services.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func getProjectID(c *gin.Context) (uuid.UUID, error) {
	projectIDStr := c.Param("projectID")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return projectID, nil
}
