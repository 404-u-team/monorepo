package handlers

import (
	"errors"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type skillsHandler struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewSkillsHandler(d *gorm.DB, l *log.Logger) skillsHandler {
	return skillsHandler{db: d, logger: l}
}

func (ch *skillsHandler) GetSkills(context *gin.Context) {
	var req dto.SkillCategoriesListRequest
	bindErr := context.ShouldBindQuery(&req)

	if bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Битый JSON"})
		return
	}

	skills, dbError := services.GetSkills(req, ch.db)

	if dbError != nil {
		if errors.Is(dbError, gorm.ErrRecordNotFound) {
			context.Status(http.StatusOK)
		} else {
			context.Status(http.StatusInternalServerError)
			ch.logger.Println("Ошибка обращения к БД:", dbError.Error())
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
	intId, err := strconv.Atoi(id)

	if id == "" || err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "id - не число"})
		return
	}

	skill, dbErr := services.GetSkillById(intId, ch.db)

	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			context.Status(http.StatusOK)
		} else {
			context.Status(http.StatusInternalServerError)
			ch.logger.Println("Ошибка обращения к БД:", dbErr.Error())
		}
		return
	}

	context.JSON(http.StatusOK, skill)
}
