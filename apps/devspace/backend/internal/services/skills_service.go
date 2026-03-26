package services

import (
	"errors"
	"math"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

func GetSkills(query dto.SkillCategoriesListQuery, db *gorm.DB) ([]dto.SkillCategoryResponse, error) {
	result := db.Model(&models.SkillCategory{})
	if query.ParentId == nil {
		result = result.Where("parent_id IS null")
	} else {
		result = result.Where("parent_id = ?", query.ParentId)
	}

	if query.Search != nil {
		result = result.Where("name ILIKE ?", *query.Search+"%")
	}

	if query.Limit != nil {
		result = result.Limit(int(*query.Limit))
	}

	var skills []models.SkillCategory

	res := result.Find(&skills)

	if res.Error != nil {
		return nil, res.Error
	}

	skillCategoryResponse := repository.BuildSkillTree(skills)
	return skillCategoryResponse, nil
}

func CutIntoPages(skills []models.SkillCategory, pages int) ([][]models.SkillCategory, error) {

	if pages == 0 {
		return nil, errors.New("0 страниц")
	}

	if len(skills) == 0 {
		return [][]models.SkillCategory{}, nil
	}
	skillsPerPage := int(math.Ceil(float64(len(skills)) / float64(pages)))

	splittedByPages := make([][]models.SkillCategory, pages)

	cursor := 0
	for idx, elem := range skills {
		if (idx+1)%skillsPerPage == 0 {
			cursor++
		}
		splittedByPages[cursor] = append(splittedByPages[cursor], elem)
	}

	return splittedByPages, nil
}

func GetSkillById(id uuid.UUID, db *gorm.DB) (*models.SkillCategory, error) {
	var targetSkill models.SkillCategory
	if res := db.Table("Skill_Category").Where("id = ?", id).First(&targetSkill); res.Error != nil {
		return nil, res.Error
	}
	return &targetSkill, nil
}

func CreateSkill(name string, parentUUID *uuid.UUID, db *gorm.DB) (*models.SkillCategory, error) {
	skill := models.SkillCategory{Name: name, ParentID: parentUUID}
	res := db.Table("Skill_Category").Create(&skill)

	if res.Error != nil {
		return nil, res.Error
	}

	return &skill, nil
}

func DeleteSkill(id uuid.UUID, cascade bool, db *gorm.DB) error {
	if cascade {
		resCascade := repository.DeleteAll("parent_id", "Skill_Category", id, db)

		if resCascade != nil {
			return resCascade
		}
	}

	res := db.Table("Skill_Category").Where("id = ?", id).Delete(&models.SkillCategory{})
	return res.Error
}

func AddSkillToUser(skillId uuid.UUID, userId uuid.UUID, db *gorm.DB) (*models.UserSkill, error) {

	var row models.UserSkill
	exists := db.Model(&models.UserSkill{}).Where("user_id = ?", userId).Where("skill_id = ?", skillId).First(&row)

	userSkill := models.UserSkill{UserID: userId, SkillID: skillId}
	if exists.Error != nil {
		if errors.Is(exists.Error, gorm.ErrRecordNotFound) {
			res := db.Model(&models.UserSkill{}).Create(&userSkill)
			return &userSkill, res.Error
		} else {
			return nil, exists.Error
		}
	}

	return nil, ErrRowAlreadyExists
}

func DeleteSkillFromUser(skillId uuid.UUID, userId uuid.UUID, db *gorm.DB) error {
	res := db.Model(&models.UserSkill{}).Where("user_id = ?", userId).Where("skill_id = ?", skillId).Delete(&models.UserSkill{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrRowNotExists
	}

	return nil
}
