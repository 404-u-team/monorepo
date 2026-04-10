package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/middleware"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type skillsHandler struct {
	db *gorm.DB
}

func NewSkillsHandler(d *gorm.DB) skillsHandler {
	return skillsHandler{db: d}
}

func (ch *skillsHandler) GetSkills(context *gin.Context) {
	var query dto.SkillCategoriesListQuery
	bindErr := context.ShouldBindQuery(&query)

	if bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка query параметров"})
		return
	}

	skills, err := services.GetSkills(query, ch.db)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.Status(http.StatusOK)
		} else {
			context.Status(http.StatusInternalServerError)
			log.Println("Ошибка обращения к БД:", err.Error())
		}
		return
	}

	context.JSON(http.StatusOK, skills)
}

func (ch *skillsHandler) GetSkillByID(context *gin.Context) {
	id := context.Param("id")
	UUID, bindErr := uuid.Parse(id)

	if bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "UUID не валиден"})
		return
	}

	skill, dbErr := services.GetSkillById(UUID, ch.db)

	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
		} else {
			context.Status(http.StatusInternalServerError)
			log.Println("Ошибка обращения к БД:", dbErr.Error())
		}
		return
	}

	context.JSON(http.StatusOK, skill)
}

func (ch *skillsHandler) CreateSkill(context *gin.Context) {
	var req dto.SkillCategoryAddRequest

	bindErr := context.ShouldBindJSON(&req)
	if bindErr != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	var dbErr error
	var skill *models.SkillCategory

	skill, dbErr = services.CreateSkill(&req, ch.db)
	if dbErr != nil {
		//Навык с таким именем уже есть
		if errors.Is(dbErr, gorm.ErrDuplicatedKey) {
			context.JSON(http.StatusConflict, gin.H{"error": "Навык с таким именем уже существует"})
			return
		} else if errors.Is(dbErr, gorm.ErrForeignKeyViolated) {
			//родителя с таким uuid нет
			context.JSON(http.StatusBadRequest, gin.H{"error": "Родителя с таким uuid не существует"})
		} else {
			// косяк на стороне сервера
			context.Status(http.StatusInternalServerError)
			log.Println("Ошибка при вставке навыка в БД: ", dbErr.Error())
		}
		return
	}

	context.JSON(http.StatusCreated, skill)
}

func (ch *skillsHandler) DeleteSkill(context *gin.Context) {
	rawUUID := context.Param("id")
	if rawUUID == "" {
		context.Status(http.StatusBadRequest)
	}

	UUID := uuid.MustParse(rawUUID)

	enableCascade := strings.ToLower(context.Query("cascade")) == "true"

	deleteError := services.DeleteSkill(UUID, enableCascade, ch.db)

	if deleteError != nil {
		if errors.Is(deleteError, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "Записи с таким id не существует"})
		} else {
			context.Status(http.StatusInternalServerError)
			if enableCascade {
				log.Println("Ошибка каскадного удаления навыка: ", deleteError.Error())
			} else {
				log.Println("Ошибка удаления навыка: ", deleteError.Error())
			}
		}

		return
	}

	context.Status(http.StatusNoContent)
}

func (ch *skillsHandler) AddSkillToSelf(context *gin.Context) {
	var req dto.BaseSkillRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	// это защищенный путь, userID не может не существовать
	userID, _ := context.Get(middleware.UserIdKey)
	userSkill, dbErr := services.AddSkillToUser(req.SkillID, userID.(uuid.UUID), ch.db)

	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrForeignKeyViolated) {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Навыка с таким uuid не существует"})
		} else if errors.Is(dbErr, services.ErrRowAlreadyExists) {
			context.JSON(http.StatusConflict, gin.H{"error": dbErr.Error()})
		} else {
			context.Status(http.StatusInternalServerError)
			log.Println("Ошибка добавления навыка пользователю в БД:" + dbErr.Error())
		}
		return
	}

	context.JSON(http.StatusCreated, userSkill)
}

func (ch *skillsHandler) DeleteSelfSkill(context *gin.Context) {
	var req dto.BaseSkillRequest
	res := context.ShouldBindQuery(&req)

	if res != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	// это защищенный путь, userID не может не существовать
	userID, _ := context.Get(middleware.UserIdKey)

	dbErr := services.DeleteSkillFromUser(req.SkillID, userID.(uuid.UUID), ch.db)
	if dbErr != nil {
		if errors.Is(dbErr, services.ErrRowNotExists) {
			context.JSON(http.StatusBadRequest, gin.H{"error": "у пользователя нет данного навыка"})
		} else {
			context.Status(http.StatusInternalServerError)
			log.Println("Ошибка удаления навыка пользователя из БД: ", dbErr.Error())
		}
		return
	}
	context.Status(http.StatusNoContent)
}
