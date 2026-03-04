package handlers

import (
	"errors"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type SkillsHandler struct {
	db *gorm.DB
}

func NewSkillsHandler(d *gorm.DB) SkillsHandler {
	return SkillsHandler{db: d}
}

func (ch *SkillsHandler) GetSkills(context *gin.Context) {
	var req dto.SkillCategoriesListRequest
	bindErr := context.ShouldBindQuery(&req)

	if bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Битый JSON"})
		return
	}

	skills, dbError := services.GetSkills(req, ch.db)

	if dbError.Error != nil {
		if errors.Is(dbError, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusOK, nil)
		} else {
			context.JSON(http.StatusInternalServerError, nil)
			log.Error(dbError)
		}
		return
	}

	if req.Page != nil {
		//делим на страницы
		splittedByPages := services.CutIntoPages(skills, int(*req.Page))
		context.JSON(http.StatusOK, splittedByPages)
		return
	}

	context.JSON(http.StatusOK, skills)
}

func (ch *SkillsHandler) GetSkillByID(context *gin.Context) {
	id := context.Param("id")
	intId, err := strconv.Atoi(id)

	if id == "" || err != nil {
		context.JSON(http.StatusBadRequest, nil)
		return
	}

	skill, dbErr := services.GetSkillById(intId, ch.db)

	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusOK, nil)
		} else {
			context.JSON(http.StatusInternalServerError, nil)
		}
		return
	}

	context.JSON(http.StatusOK, skill)
}
