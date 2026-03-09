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

type projectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) *projectHandler {
	return &projectHandler{
		projectService: projectService,
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

	err = h.projectService.UpdateProjectByID(projectID, userID, &payload)
	if err != nil {
		if errors.Is(err, services.ErrUserNotLeader) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		if errors.Is(err, services.ErrProjectConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *projectHandler) DeleteProjectByID(c *gin.Context) {
	// получение projectID из параметров
	projectID, err := getProjectID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrProjectNotFound.Error()})
		return
	}

	// получение ID пользователя из контекста
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	err = h.projectService.DeleteProjectByID(projectID, userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotLeader) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		if errors.Is(err, services.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrProjectHasSlots) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func getProjectID(c *gin.Context) (uuid.UUID, error) {
	projectIDStr := c.Param("projectID")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return projectID, nil
}
