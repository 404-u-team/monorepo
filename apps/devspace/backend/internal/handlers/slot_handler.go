package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type slotHandler struct {
	slotService    services.SlotService
	projectService services.ProjectService
}

func NewSlotHandler(slotService services.SlotService, projectService services.ProjectService) *slotHandler {
	return &slotHandler{
		slotService:    slotService,
		projectService: projectService,
	}
}

func (h *slotHandler) GetSlots(c *gin.Context) {
	projectID, err := getProjectID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrProjectNotFound.Error()})
		return
	}

	projectSlots, err := h.slotService.GetSlots(projectID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, projectSlots)
}

func (h *slotHandler) CreateSlot(c *gin.Context) {
	var payload dto.CreateSlotRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	projectID, err := getProjectID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrProjectNotFound.Error()})
		return
	}

	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := h.projectService.IsUserProjectLeader(projectID, userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if !isUserProjectLeader {
		c.Status(http.StatusForbidden)
		return
	}

	err = h.slotService.CreateSlot(projectID, &payload)
	if err != nil {
		if errors.Is(err, services.ErrSlotConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
