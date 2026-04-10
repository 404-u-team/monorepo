package services

import (
	"errors"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

// TODO: пиздец. Не работает поиск по любым родителям (только корень или только ребенок определенного скила)
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

	if query.StartAt != nil {
		result = result.Offset(*query.StartAt)
	}

	if query.Limit != nil {
		result = result.Limit(*query.Limit)
	}

	var skills []models.SkillCategory

	res := result.Find(&skills)

	if res.Error != nil {
		return nil, res.Error
	}

	skillCategoryResponse := repository.BuildSkillTree(skills)

	return skillCategoryResponse, nil
}

func GetSkillById(id uuid.UUID, db *gorm.DB) (*models.SkillCategory, error) {
	var targetSkill models.SkillCategory
	if res := db.Table("Skill_Category").Where("id = ?", id).First(&targetSkill); res.Error != nil {
		return nil, res.Error
	}
	return &targetSkill, nil
}

func CreateSkill(req *dto.SkillCategoryAddRequest, db *gorm.DB) (*models.SkillCategory, error) {
	skill := models.SkillCategory{Name: req.Name, ParentID: req.ParentID, Icon: req.Icon, Color: req.Color}
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
