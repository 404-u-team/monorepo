package repository

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IdeaRepository interface {
	UpdateIdeaByID(ideaID uuid.UUID, updateRequest *dto.UpdateIdeaRequest) (int, error)
	IsUserIdeaAuthor(ideaID, userID uuid.UUID) (bool, error)
	GetIdeaByID(id uuid.UUID) (*models.Idea, error)
	GetIdeas(query *dto.GetIdeasRequest, userID uuid.UUID) ([]models.Idea, int64, error)
	ToggleFavorite(ideaID, userID uuid.UUID) (bool, error)
}

type ideaRepository struct {
	conn *gorm.DB
}

func NewIdeaRepository(conn *gorm.DB) IdeaRepository {
	return &ideaRepository{conn: conn}
}

// обновить данные идеи, возвращает кол-во измененных строк и ошибку
func (r *ideaRepository) UpdateIdeaByID(ideaID uuid.UUID, updateRequest *dto.UpdateIdeaRequest) (int, error) {
	updates := map[string]interface{}{}

	if updateRequest.Title != nil {
		updates["title"] = *updateRequest.Title
	}

	if updateRequest.Description != nil {
		updates["description"] = *updateRequest.Description
	}

	result := r.conn.Model(&models.Idea{}).Where("id = ?", ideaID).Updates(updates)
	if result.Error != nil {
		log.Println("Ошибка при обновлении идеи по ID: ", result.Error)
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
		log.Println("Ошибка при проверке является ли пользователь владельцем идеи: ", result.Error)
		return false, result.Error
	}
	return count == 1, nil
}

func (r *ideaRepository) GetIdeaByID(id uuid.UUID) (*models.Idea, error) {
	var idea models.Idea

	res := r.conn.Model(&models.Idea{}).Where("id = ?", id).First(&idea)

	if res.Error != nil {
		log.Println("Ошибка при получении идеи по ID: ", res.Error)
		return nil, res.Error
	}
	return &idea, nil
}

func (r *ideaRepository) GetIdeas(query *dto.GetIdeasRequest, userID uuid.UUID) ([]models.Idea, int64, error) {
	result := r.conn.Model(&models.Idea{})

	if query.Search != nil {
		searchInner := "%" + (*query.Search) + "%"
		result = result.Where("(title ILIKE ? OR description ILIKE ? OR content ILIKE ? OR category ILIKE ?)", searchInner, searchInner, searchInner, searchInner)
	}

	if query.AuthorId != nil {
		result = result.Where("author_id = ?", *query.AuthorId)
	}

	if query.IsFavorite {
		result = result.Joins(`JOIN "User_Favorite" ON "User_Favorite".idea_id = "Idea".id`).
			Where(`"User_Favorite".user_id = ?`, userID)
	}

	var total int64
	if err := result.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.StartAt != nil {
		result = result.Offset(int(*query.StartAt))
	}

	if query.Limit != nil {
		result = result.Limit(int(*query.Limit))
	}

	var ideas []models.Idea
	if err := result.Find(&ideas).Error; err != nil {
		return nil, 0, err
	}

	return ideas, total, nil
}

func (r *ideaRepository) ToggleFavorite(ideaID, userID uuid.UUID) (bool, error) {
	// создаем транзакцию, у которой есть два исхода
	// 		1. Если есть строчка с idea_id = true – удаляем строчку
	// 		2. Если не нашли ничего, то создаем новую строчку
	var isFavorite bool // произошло удаление или создание
	err := r.conn.Raw(`
		WITH deleted AS (
			DELETE FROM "User_Favorite_Idea"
			WHERE user_id = ? AND idea_id = ?
			RETURNING FALSE AS action
		),
		inserted AS (
			INSERT INTO "User_Favorite_Idea" (user_id, idea_id)
			SELECT ?, ?
			WHERE NOT EXISTS (SELECT 1 FROM deleted)
			RETURNING TRUE AS action
		)
		SELECT action FROM deleted
		UNION ALL
		SELECT action FROM inserted
	`, userID, ideaID, userID, ideaID).Scan(&isFavorite).Error

	if err != nil {
		return false, err
	}

	return isFavorite, nil
}
