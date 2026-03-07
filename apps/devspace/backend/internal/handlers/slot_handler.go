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
		if errors.Is(err, services.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *slotHandler) UpdateSlotByID(c *gin.Context) {
	// получение slotID из параметров
	slotID, err := getSlotID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrSlotNotFound.Error()})
		return
	}

	// получение projectID из параметров
	projectID, err := getProjectID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrProjectNotFound.Error()})
		return
	}

	// получение projectID из контекста
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
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

	// получение payload
	var payload dto.UpdateSlotRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	if payload.SkillCategoryID == nil && payload.Title == nil && payload.Description == nil && payload.Status == nil {
		log.Println("Все поля пустые, нечего изменять")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля пустые, нечего изменять"})
		return
	}

	err = h.slotService.UpdateSlotByID(slotID, &payload)
	if err != nil {
		if errors.Is(err, services.ErrSlotConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		if errors.Is(err, services.ErrSlotNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *slotHandler) DeleteSlotByID(c *gin.Context) {
	// получение slotID из параметров
	slotID, err := getSlotID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrSlotNotFound.Error()})
		return
	}

	// получение projectID из параметров
	projectID, err := getProjectID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrProjectNotFound.Error()})
		return
	}

	// получение projectID из контекста
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
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

	err = h.slotService.DeleteSlotByID(slotID)
	if err != nil {
		if errors.Is(err, services.ErrSlotNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func getSlotID(c *gin.Context) (uuid.UUID, error) {
	slotIDStr := c.Param("slotID")

	slotID, err := uuid.Parse(slotIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return slotID, nil
}
