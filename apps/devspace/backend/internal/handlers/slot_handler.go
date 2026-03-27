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
	slotService services.SlotService
}

func NewSlotHandler(slotService services.SlotService) *slotHandler {
	return &slotHandler{
		slotService: slotService,
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

	slot, err := h.slotService.CreateSlot(projectID, userID, &payload)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSlotSkills) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrUserNotLeader) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrSlotConflict) {
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

	c.JSON(http.StatusCreated, slot)
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

	// получение userID из контекста
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// получение payload
	var payload dto.UpdateSlotRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	slot, err := h.slotService.UpdateSlotByID(slotID, projectID, userID, &payload)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSlotSkills) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrEmptyPayload) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrSlotConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrSlotNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrUserNotLeader) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, slot)
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

	// получение userID из контекста
	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	err = h.slotService.DeleteSlotByID(slotID, projectID, userID)
	if err != nil {
		if errors.Is(err, services.ErrSlotNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrUserNotLeader) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
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
