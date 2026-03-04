package services

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
	"math"
)

func GetSkills(req dto.SkillCategoriesListRequest, db *gorm.DB) ([]models.SkillCategory, error) {
	query := db.Table("skill_categories")
	if req.ParentId == nil {
		query = query.Where("parent_id IS null")
	} else {
		query = query.Where("parent_id = ?", req.ParentId)
	}

	if req.Search != nil {
		query = query.Where("name LIKE ?%", req.Search)
	}

	if req.Limit != nil {
		query = query.Limit(int(*req.Limit))
	}

	var skills []models.SkillCategory

	res := query.Find(&skills)

	if res.Error != nil {
		return nil, res.Error
	}

	return skills, nil
}

func CutIntoPages(skills []models.SkillCategory, pages int) [][]models.SkillCategory {
	skillsPerPage := int(math.Ceil(float64(len(skills)) / float64(pages)))

	splittedByPages := make([][]models.SkillCategory, pages)

	cursor := 0
	for idx, elem := range skills {
		if (idx+1)%skillsPerPage == 0 {
			cursor++
		}
		splittedByPages[cursor] = append(splittedByPages[cursor], elem)
	}

	return splittedByPages
}

func GetSkillById(id int, db *gorm.DB) (*models.SkillCategory, error) {
	var targetSkill models.SkillCategory
	if res := db.Table("Skill_Category").Where("id = ?", id).First(&targetSkill); res.Error != nil {
		return nil, res.Error
	}
	return &targetSkill, nil
}
