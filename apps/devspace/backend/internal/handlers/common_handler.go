package handlers

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
	"gorm.io/gorm"
)

type CommonHandler struct {
	db *gorm.DB
}

func NewCommonHandler(d *gorm.DB) CommonHandler {
	return CommonHandler{db: d}
}

func (ch *CommonHandler) GetSkills(context *gin.Context) {
	var req dto.SkillCategoriesListRequest
	bindErr := context.ShouldBindJSON(&req)

	if bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Битый JSON"})
		return
	}

	queryParts := utils.MakeSlice(5, "SELECT * FROM")
	params := make([]any, 0, 4)

	if req.ParentId != nil && req.Search != nil {
		queryParts = append(queryParts, "WHERE parent_id = ?", "AND", "LIKE ?%")
		params = append(params, req.ParentId, req.Search)
	} else {
		if req.ParentId == nil {
			queryParts = append(queryParts, "WHERE parent_id IS null")
		} else {
			queryParts = append(queryParts, "WHERE parent_id = ?")
			params = append(params, req.ParentId)
		}

		if req.Search != nil {
			queryParts = append(queryParts, "AND", "name LIKE ?%")
			params = append(params, req.Search)
		}
	}

	if req.Limit != nil {
		queryParts = append(queryParts, "LIMIT ?")
		params = append(params, req.Limit)
	}

	var skills []models.SkillCategory

	res := ch.db.Raw(strings.Join(queryParts, " "), params...).Scan(&skills)

	if res.Error != nil {
		log.Error(res.Error)
		context.JSON(http.StatusInternalServerError, nil)
		return
	}

	if req.Page != nil {
		//делим на страницы
		skillsPerPage := int(math.Ceil(float64(len(skills)) / float64(*req.Page)))

		splittedByPages := make([][]models.SkillCategory, *req.Page)

		cursor := 0
		for idx, elem := range skills {
			if (idx+1)%skillsPerPage == 0 {
				cursor++
			}
			splittedByPages[cursor] = append(splittedByPages[cursor], elem)
		}

		context.JSON(http.StatusOK, splittedByPages)
		return
	}

	context.JSON(http.StatusOK, skills)
}

func (ch *CommonHandler) GetSkillByID(context *gin.Context) {
	id := context.Param("id")
	intId, err := strconv.Atoi(id)

	if id == "" || err != nil {
		context.JSON(http.StatusBadRequest, nil)
		return
	}

	var targetSkill models.SkillCategory

	if res := ch.db.Table("Skill_Category").Where("id = ?", intId).First(&targetSkill); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, nil)
		} else {
			context.JSON(http.StatusInternalServerError, nil)
		}
		return
	}

	context.JSON(http.StatusOK, targetSkill)
}
