package services

import (
	"errors"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IdeaService interface {
	UpdateIdeaByID(ideaID, userID uuid.UUID, updateRequest *dto.UpdateIdeaRequest) (*models.Idea, error)
}

type ideaService struct {
	repo repository.IdeaRepository
}

func NewIdeaService(repo repository.IdeaRepository) IdeaService {
	return &ideaService{repo: repo}
}

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
	idea := models.Idea{AuthorID: authorId, Title: req.Title, Description: req.Description}

	if req.Content != nil {
		idea.Content = req.Content
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

func GetIdeaByID(id uuid.UUID, db *gorm.DB) (*models.Idea, error) {
	var idea models.Idea

	res := db.Where("id = ?", id).First(&idea)

	if res.Error != nil {
		return nil, res.Error
	}
	return &idea, nil
}

func (s *ideaService) UpdateIdeaByID(ideaID, userID uuid.UUID, updateRequest *dto.UpdateIdeaRequest) (*models.Idea, error) {
	// является ли пользователь владельцем данной идеи
	isUserIdeaAuthor, err := s.repo.IsUserIdeaAuthor(ideaID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isUserIdeaAuthor {
		return nil, ErrUserNotAuthor
	}

	// обновление идеи по ID
	rowsAffected, err := s.repo.UpdateIdeaByID(ideaID, updateRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrIdeaConflict
		}
		return nil, ErrInternal
	}

	if rowsAffected == 0 {
		return nil, ErrIdeaNotFound
	}

	// получаем обновленную идею для возвращения
	idea, err := s.repo.GetIdeaByID(ideaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIdeaNotFound
		}
		return nil, ErrInternal
	}

	return idea, nil
}
