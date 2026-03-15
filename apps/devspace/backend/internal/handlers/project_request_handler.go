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

type projectRequestHandler struct {
	projectRequestService services.ProjectRequestService
}

func NewProjectRequestHandler(projectRequestService services.ProjectRequestService) *projectRequestHandler {
	return &projectRequestHandler{
		projectRequestService: projectRequestService,
	}
}

func (h *projectRequestHandler) CreateProjectRequestApply(c *gin.Context) {
	var payload dto.CreateProjectRequestApplyRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	projectIDStr := c.Param("projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	slotIDStr := c.Param("slotID")
	slotID, err := uuid.Parse(slotIDStr)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	projectRequest, err := h.projectRequestService.CreateProjectRequestApply(&payload, slotID, userID, projectID)
	if err != nil {
		if errors.Is(err, services.ErrSlotNotFound) || errors.Is(err, services.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrSlotIsClosed) || errors.Is(err, services.ErrUserLeader) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrProjectRequestConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, projectRequest)
}

func (h *projectRequestHandler) CreateProjectRequestInvite(c *gin.Context) {
	var payload dto.CreateProjectRequestInviteRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	projectIDStr := c.Param("projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	slotIDStr := c.Param("slotID")
	slotID, err := uuid.Parse(slotIDStr)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	projectRequest, err := h.projectRequestService.CreateProjectRequestInvite(&payload, slotID, userID, projectID)
	if err != nil {
		if errors.Is(err, services.ErrSlotNotFound) || errors.Is(err, services.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrSlotIsClosed) || errors.Is(err, services.ErrCantInviteYourself) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrProjectRequestConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrUserNotLeader) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, projectRequest)
}
