package repository

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IdeaRepository interface {
	UpdateIdeaByID(ideaID uuid.UUID, updateRequest *dto.UpdateIdeaRequest) (int, error)
	IsUserIdeaAuthor(ideaID, userID uuid.UUID) (bool, error)
	GetIdeaByID(id uuid.UUID) (*models.Idea, error)
}

type ideaRepository struct {
	conn *gorm.DB
}

func NewIdeaRepository(conn *gorm.DB) IdeaRepository {
	return &ideaRepository{conn: conn}
}

// обновить данные идеи, возвращает кол-во измененных строк и ошибку
func (r *ideaRepository) UpdateIdeaByID(ideaID uuid.UUID, updateRequest *dto.UpdateIdeaRequest) (int, error) {
	updates := map[string]string{}

	if updateRequest.Title != nil {
		updates["title"] = *updateRequest.Title
	}

	if updateRequest.Description != nil {
		updates["description"] = *updateRequest.Description
	}

	result := r.conn.Model(&models.Idea{}).Where("id = ?", ideaID).Updates(updates)
	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, nil
	}

	return int(result.RowsAffected), nil
}

func (r *ideaRepository) IsUserIdeaAuthor(ideaID, userID uuid.UUID) (bool, error) {
	var count int64
	result := r.conn.Model(&models.Idea{}).Where("id = ?", ideaID).Where("author_id = ?", userID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count == 1, nil
}

func (r *ideaRepository) GetIdeaByID(id uuid.UUID) (*models.Idea, error) {
	var idea models.Idea

	res := r.conn.Model(&models.Idea{}).Where("id = ?", id).First(&idea)

	if res.Error != nil {
		return nil, res.Error
	}
	return &idea, nil
}
