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
	GetIdeaByID(ideaID, userID uuid.UUID) (*dto.GetIdeaResponse, error)
	GetIdeaByIDIncr(ideaID, userID uuid.UUID) (*dto.GetIdeaResponse, error)
	GetIdeas(query *dto.GetIdeasRequest, userID uuid.UUID) ([]dto.IdeaBlock, int64, error)
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

func (r *ideaRepository) GetIdeaByID(ideaID, userID uuid.UUID) (*dto.GetIdeaResponse, error) {
	var ideaResponse dto.GetIdeaResponse
	result := r.conn.Table("Idea").Select(`id, author_id, author_id = ? AS is_author, COALESCE("User_Favorite_Idea".idea_id IS NOT NULL, false) AS is_favorite, 
		title, description, views_count, favorites_count, category, created_at, updated_at`, userID).
		Joins(`LEFT JOIN "User_Favorite_Idea" ON "User_Favorite_Idea".idea_id = "Idea".id AND "User_Favorite_Idea".user_id = ?`, userID).
		Where("id = ?", ideaID)

	if err := result.First(&ideaResponse).Error; err != nil {
		log.Println("Ошибка при получении идеи по ID: ", err)
		return nil, err
	}
	return &ideaResponse, nil
}

// такая же версия как и GetIdeaByID, но с инкрементацией счетчика views
func (r *ideaRepository) GetIdeaByIDIncr(ideaID, userID uuid.UUID) (*dto.GetIdeaResponse, error) {
	var ideaResponse dto.GetIdeaResponse
	result := r.conn.Raw(`
		UPDATE "Idea" i
		SET views_count = views_count + 1
		FROM (
			SELECT 
				i2.id,
				i2.author_id = ? AS is_author,
				COALESCE(uf.idea_id IS NOT NULL, false) AS is_favorite
			FROM "Idea" i2
			LEFT JOIN "User_Favorite_Idea" uf 
				ON uf.idea_id = i2.id AND uf.user_id = ?
			WHERE i2.id = ?
		) AS sub
		WHERE i.id = sub.id
		RETURNING 
			i.id,
			i.author_id,
			sub.is_author,
			sub.is_favorite,
			i.title,
			i.description,
			i.views_count,
			i.favorites_count,
			i.category,
			i.created_at,
			i.updated_at
	`, userID, userID, ideaID)

	if err := result.First(&ideaResponse).Error; err != nil {
		log.Println("Ошибка при получении идеи по ID: ", err)
		return nil, err
	}
	return &ideaResponse, nil
}

func (r *ideaRepository) GetIdeas(query *dto.GetIdeasRequest, userID uuid.UUID) ([]dto.IdeaBlock, int64, error) {
	result := r.conn.Table("Idea").
		Select(`id, author_id, author_id = ? AS is_author, COALESCE("User_Favorite_Idea".idea_id IS NOT NULL, false) AS is_favorite, 
		title, description, views_count, favorites_count, category, created_at, updated_at`, userID).
		Joins(`LEFT JOIN "User_Favorite_Idea" ON "User_Favorite_Idea".idea_id = "Idea".id AND "User_Favorite_Idea".user_id = ?`, userID)

	if query.Search != nil {
		searchInner := "%" + (*query.Search) + "%"
		result = result.Where("(title ILIKE ? OR description ILIKE ? OR content ILIKE ? OR category ILIKE ?)", searchInner, searchInner, searchInner, searchInner)
	}

	if query.AuthorId != nil {
		result = result.Where("author_id = ?", *query.AuthorId)
	}

	if query.IsFavorite {
		result = result.Where(`"User_Favorite_Idea".idea_id IS NOT NULL`)
	}

	if query.Views != nil {
		if *query.Views == "asc" {
			result = result.Order("views_count ASC")
		} else {
			result = result.Order("views_count DESC")
		}
	}

	if query.Favorites != nil {
		if *query.Favorites == "asc" {
			result = result.Order("favorites_count ASC")
		} else {
			result = result.Order("favorites_count DESC")
		}
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

	var ideasBlock []dto.IdeaBlock
	if err := result.Find(&ideasBlock).Error; err != nil {
		log.Println("Произошла ошибка при получении списка идей: ", err)
		return nil, 0, err
	}

	return ideasBlock, total, nil
}

func (r *ideaRepository) ToggleFavorite(ideaID, userID uuid.UUID) (bool, error) {
	// создаем транзакцию, у которой есть два исхода
	// 		1. Если есть строчка с idea_id = true – удаляем строчку
	// 		2. Если не нашли ничего, то создаем новую строчку
	var isFavorite bool // произошло удаление или создание
	err := r.conn.Raw(`
		WITH
		deleted AS (
			DELETE FROM "User_Favorite_Idea"
			WHERE user_id = ? AND idea_id = ?
			RETURNING idea_id
		),
		decremented AS (
			UPDATE "Idea"
			SET favorites_count = favorites_count - 1
			WHERE id IN (SELECT idea_id FROM deleted)
			RETURNING FALSE AS action
		),
		inserted AS (
			INSERT INTO "User_Favorite_Idea" (user_id, idea_id)
			SELECT ?, ?
			WHERE NOT EXISTS (SELECT 1 FROM deleted)
			RETURNING idea_id
		),
		incremented AS (
			UPDATE "Idea"
			SET favorites_count = favorites_count + 1
			WHERE id IN (SELECT idea_id FROM inserted)
			RETURNING TRUE AS action
		)
		SELECT action FROM decremented
		UNION ALL
		SELECT action FROM incremented;	
	`, userID, ideaID, userID, ideaID).Scan(&isFavorite).Error

	if err != nil {
		log.Println("Произошла ошибка при toggle favorite idea: ", err)
		return false, err
	}

	return isFavorite, nil
}
