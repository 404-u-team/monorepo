package repository

import (
	"errors"
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillRepository interface {
	GetSkillByID(skillID uuid.UUID) (*models.SkillCategory, error)
}

type skillRepository struct {
	conn *gorm.DB
}

func NewSkillRepository(conn *gorm.DB) SkillRepository {
	return &skillRepository{conn: conn}
}

func (r *skillRepository) GetSkillByID(skillID uuid.UUID) (*models.SkillCategory, error) {
	var skill models.SkillCategory
	result := r.conn.Model(&models.SkillCategory{}).Where("id = ?", skillID).First(&skill)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Println("Ошибка при получении скилла по ID: ", result.Error)
		return nil, result.Error
	}

	return &skill, nil
}
