package repository

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SlotRepository interface {
	GetSlots(projectID uuid.UUID) ([]models.ProjectSlot, error)
	CreateSlot(projectID uuid.UUID, slot *models.ProjectSlot) error
	UpdateSlotByID(slotID uuid.UUID, updateRequest *dto.UpdateSlotRequest) (int, error)
	DeleteSlotByID(slotID uuid.UUID) (int, error)
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

// обновить слот для определенного проекта. Возвращает количество изменных строк и ошибку
func (r *slotRepository) UpdateSlotByID(slotID uuid.UUID, updateRequest *dto.UpdateSlotRequest) (int, error) {
	updates := map[string]string{}

	if updateRequest.SkillCategoryID != nil {
		updates["skill_category_id"] = (*updateRequest.SkillCategoryID).String()
	}

	if updateRequest.Title != nil {
		updates["title"] = *updateRequest.Title
	}

	if updateRequest.Description != nil {
		updates["description"] = *updateRequest.Description
	}

	if updateRequest.Status != nil {
		updates["status"] = *updateRequest.Status
	}

	result := r.conn.Model(&models.ProjectSlot{}).Where("id = ?", slotID).Updates(updates)
	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, nil
	}

	return int(result.RowsAffected), nil
}

// возвращает количество удаленных проектов и ошибку
func (r *slotRepository) DeleteSlotByID(slotID uuid.UUID) (int, error) {
	result := r.conn.Delete(&models.ProjectSlot{}, "id = ?", slotID)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}
