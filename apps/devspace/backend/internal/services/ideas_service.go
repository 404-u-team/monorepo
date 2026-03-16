package services

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetIdeasList(req dto.GetListIdeasRequest, db *gorm.DB) ([]models.Idea, error) {
	query := db.Model(&models.Idea{})
	var ideas []models.Idea

	if req.Search != nil {
		query = query.Where("title ILIKE ?", *req.Search+"%")
	}

	if req.AuthorId != nil {
		query = query.Where("author_id = ?", *req.AuthorId)
	}

	if req.StartAt != nil {
		query = query.Offset(int(*req.StartAt))
	}

	if req.Limit != nil {
		query = query.Limit(int(*req.Limit))
	}

	res := query.Find(&ideas)

	if res.Error != nil {
		return nil, res.Error
	}
	return ideas, nil
}

func CreateIdea(req dto.CreateIdeaRequest, authorId uuid.UUID, db *gorm.DB) (*models.Idea, error) {
	idea := models.Idea{AuthorID: authorId, Title: req.Title}

	if req.Description != nil {
		idea.Description = *req.Description
	}

	if req.Category != nil {
		idea.Category = *req.Category
	}
	// в ходе create gorm скорректирует нужные поля у сущности, вроде id
	res := db.Create(&idea)

	if res.Error != nil {
		return nil, res.Error
	}

	return &idea, nil
}
