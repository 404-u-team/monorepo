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
		c.Status(http.StatusBadRequest)
		return
	}

	ideasResponse, err := ih.ideaService.GetIdeas(&query, ih.config, c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ideasResponse)
}

func (ih *ideaHandler) AddIdea(ctx *gin.Context) {
	var req dto.CreateIdeaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	userID, _ := ctx.Get("userID")

	// разыменовываем any, т.к там 100% uuid
	idea, err := services.CreateIdea(req, userID.(uuid.UUID), ih.db)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.JSON(http.StatusConflict, gin.H{"code": http.StatusConflict, "error": "Идея с таким названием уже существует"})
		} else {
			ctx.Status(http.StatusInternalServerError)
			log.Println("Ошибка записи идеи в БД: " + err.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, idea)
}

func (ih *ideaHandler) GetIdeaByID(ctx *gin.Context) {
	id := ctx.Param("id")

	converted, parseError := uuid.Parse(id)

	if parseError != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	idea, dbErr := services.GetIdeaByID(converted, ih.db)
	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
			log.Println("Ошибка получения идеи из БД по uuid: " + dbErr.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, idea)
}

func (h *ideaHandler) UpdateIdeaByID(c *gin.Context) {
	// получение ideaID из параметров
	ideaIDStr := c.Param("ideaID")

	ideaID, err := uuid.Parse(ideaIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrIdeaNotFound.Error()})
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

	idea, err := h.ideaService.UpdateIdeaByID(ideaID, userID, &payload)
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

	c.JSON(http.StatusOK, idea)
}
func (ih *ideaHandler) DeleteIdeaByID(ctx *gin.Context) {
	rawID := ctx.Param("id")
	ideaID, parseErr := uuid.Parse(rawID)

	if parseErr != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	//это защищенный путь, ID 100% существует
	userID := ctx.MustGet("userID").(uuid.UUID)

	canDelete, dbErr := services.CheckRightsOnIdea(ideaID, userID, ih.db)

	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Запись с таким id не существует"})
			return
		} else {
			ctx.Status(http.StatusInternalServerError)
			log.Println("Ошибка получения прав пользователя на удаление записи: " + dbErr.Error())
		}
		return
	}

	if !canDelete {
		ctx.Status(http.StatusForbidden)
		return
	}

	dbErr = services.DeleteIdeaByID(ideaID, ih.db)
	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "идеи с таким ID не существует"})
		} else {
			ctx.Status(http.StatusInternalServerError)
			log.Println("Ошибка удаления записи: " + dbErr.Error())
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ih *ideaHandler) ToggleFavorite(c *gin.Context) {
	ideaIDStr := c.Param("projectID")
	ideaID, err := uuid.Parse(ideaIDStr)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	isFavorite, err := ih.ideaService.ToggleFavorite(ideaID, userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	toggleFavoriteResponse := dto.ToggleFavoriteResponse{IsFavorite: isFavorite}
	c.JSON(http.StatusOK, toggleFavoriteResponse)
}
