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
		if errors.Is(err, services.ErrSlotIsClosed) || errors.Is(err, services.ErrCantInviteYourself) || errors.Is(err, services.ErrBadRequest) {
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

func (h *projectRequestHandler) AcceptProjectRequest(c *gin.Context) {
	// получение projectRequest из параметров
	projectRequestIDStr := c.Param("projectRequestID")
	projectRequestID, err := uuid.Parse(projectRequestIDStr)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// получение userID
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// одобрить заявку
	projectRequest, err := h.projectRequestService.UpdateProjectRequest(projectRequestID, userID, "accepted")
	if err != nil {
		if errors.Is(err, services.ErrUserNotLeader) || errors.Is(err, services.ErrProjectRequestDontBelongToUser) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrProjectRequestNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrProjectRequestNotPending) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, projectRequest)
}

func (h *projectRequestHandler) RejectProjectRequest(c *gin.Context) {
	// получение projectRequest из параметров
	projectRequestIDStr := c.Param("projectRequestID")
	projectRequestID, err := uuid.Parse(projectRequestIDStr)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// получение userID
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// одобрить заявку
	projectRequest, err := h.projectRequestService.UpdateProjectRequest(projectRequestID, userID, "rejected")
	if err != nil {
		if errors.Is(err, services.ErrUserNotLeader) || errors.Is(err, services.ErrProjectRequestDontBelongToUser) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrProjectRequestNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrProjectRequestNotPending) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, projectRequest)
}
