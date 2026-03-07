package handlers

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"

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
	var req dto.SkillCategoriesListRequest
	bindErr := context.ShouldBindQuery(&req)

	if bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Битый JSON"})
		return
	}

	if req.ParentId != nil {
		_, parceErr := uuid.Parse(*req.ParentId)
		if parceErr != nil {
			context.Status(http.StatusNotFound)
			return
		}
	}

	skills, dbError := services.GetSkills(req, ch.db)

	if dbError != nil {
		if errors.Is(dbError, gorm.ErrRecordNotFound) {
			context.Status(http.StatusOK)
		} else {
			context.Status(http.StatusInternalServerError)
			log.Println("Ошибка обращения к БД:", dbError.Error())
		}
		return
	}

	if req.Page != nil {
		//делим на страницы
		splittedByPages, err := services.CutIntoPages(skills, int(*req.Page))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		context.JSON(http.StatusOK, splittedByPages)
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
			context.Status(http.StatusOK)
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

	if context.Param("id") == "" {
		dbErr = services.CreateSkill(req.Name, nil, ch.db)
	} else {
		UUID, parceErr := uuid.Parse(context.Param("id"))
		if parceErr != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Битый uuid"})
			return
		}

		dbErr = services.CreateSkill(req.Name, &UUID, ch.db)
	}
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

	context.Status(http.StatusCreated)
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

	context.Status(http.StatusOK)
}
