package repository

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SlotRepository interface {
	GetSlots(projectID uuid.UUID) ([]models.ProjectSlot, error)
	CreateSlot(projectID uuid.UUID, slot *models.ProjectSlot) error
}

type slotRepository struct {
	conn *gorm.DB
}

func NewSlotRepository(conn *gorm.DB) SlotRepository {
	return &slotRepository{conn: conn}
}

func (r *slotRepository) GetSlots(projectID uuid.UUID) ([]models.ProjectSlot, error) {
	var slots []models.ProjectSlot
	result := r.conn.Model(&models.ProjectSlot{}).Where("project_id = ?", projectID).Find(&slots)
	if result.Error != nil {
		log.Println("Ошибка при получении списка слотов проекта: ", result.Error)
		return nil, result.Error
	}

	return slots, nil
}

func (r *slotRepository) CreateSlot(projectID uuid.UUID, slot *models.ProjectSlot) error {
	var count int64
	if err := r.conn.Model(&models.Project{}).
		Where("id = ?", projectID).
		Count(&count).Error; err != nil {
		log.Println("Ошибка при проверке наличия проекта: ", err)
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	result := r.conn.Create(slot)
	if result.Error != nil {
		log.Println("Ошибка при создании слота для проекта: ", result.Error)
		return result.Error
	}

	return nil
}
