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
	"gorm.io/gorm"
)

type ideaHandler struct {
	ideaService services.IdeaService
	db          *gorm.DB
	config      *config.Config
}

func NewIdeaHandler(ideaService services.IdeaService, db *gorm.DB, config *config.Config) ideaHandler {
	return ideaHandler{ideaService: ideaService, db: db, config: config}
}

func (ih *ideaHandler) GetIdeas(c *gin.Context) {
	var query dto.GetIdeasRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Println("Ошибка ShouldBindQuery: ", err)
		c.Status(http.StatusBadRequest)
		return
	}

	ideasResponse, err := ih.ideaService.GetIdeas(&query, ih.config, c)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ideasResponse)
}

func (ih *ideaHandler) CreateIdea(c *gin.Context) {
	var req dto.CreateIdeaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// разыменовываем any, т.к там 100% uuid
	idea, err := ih.ideaService.CreateIdea(&req, userID)
	if err != nil {
		if errors.Is(err, services.ErrIdeaConflict) {
			c.Status(http.StatusConflict)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, idea)
}

func (ih *ideaHandler) GetIdeaByID(c *gin.Context) {
	ideaIDStr := c.Param("ideaID")
	ideaID, parseError := uuid.Parse(ideaIDStr)
	if parseError != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	idea, err := ih.ideaService.GetIdeaByID(ideaID, ih.config, c)
	if err != nil {
		if errors.Is(err, services.ErrIdeaNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, idea)
}

func (h *ideaHandler) UpdateIdeaByID(c *gin.Context) {
	// получение ideaID из параметров
	ideaIDStr := c.Param("ideaID")

	ideaID, err := uuid.Parse(ideaIDStr)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// получение payload
	var payload dto.UpdateIdeaRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Ошибка при парсинге: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	if payload.Title == nil && payload.Description == nil {
		log.Println("Все поля пустые, нечего изменять")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля пустые, нечего изменять"})
		return
	}

	ideaResponse, err := h.ideaService.UpdateIdeaByID(ideaID, userID, &payload)
	if err != nil {
		if errors.Is(err, services.ErrUserNotAuthor) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrIdeaConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrIdeaNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ideaResponse)
}
func (ih *ideaHandler) DeleteIdeaByID(c *gin.Context) {
	ideaIDStr := c.Param("ideaID")

	ideaID, err := uuid.Parse(ideaIDStr)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//это защищенный путь, ID 100% существует
	userID := c.MustGet("userID").(uuid.UUID)

	canDelete, dbErr := services.CheckRightsOnIdea(ideaID, userID, ih.db)

	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Запись с таким id не существует"})
			return
		} else {
			c.Status(http.StatusInternalServerError)
			log.Println("Ошибка получения прав пользователя на удаление записи: " + dbErr.Error())
		}
		return
	}

	if !canDelete {
		c.Status(http.StatusForbidden)
		return
	}

	dbErr = services.DeleteIdeaByID(ideaID, ih.db)
	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "идеи с таким ID не существует"})
		} else {
			c.Status(http.StatusInternalServerError)
			log.Println("Ошибка удаления записи: " + dbErr.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (ih *ideaHandler) ToggleFavorite(c *gin.Context) {
	ideaIDStr := c.Param("ideaID")
	ideaID, err := uuid.Parse(ideaIDStr)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	toggleFavoriteResponse, err := ih.ideaService.ToggleFavorite(ideaID, userID)
	if err != nil {
		if errors.Is(err, services.ErrIdeaNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, toggleFavoriteResponse)
}
