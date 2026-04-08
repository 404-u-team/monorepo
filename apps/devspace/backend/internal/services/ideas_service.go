package services

import (
	"errors"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/middleware"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IdeaService interface {
	UpdateIdeaByID(ideaID, userID uuid.UUID, updateRequest *dto.UpdateIdeaRequest) (*dto.GetIdeaResponse, error)
	GetIdeas(query *dto.GetIdeasRequest, config *config.Config, c *gin.Context) (*dto.GetIdeasResponse, error)
	GetIdeaByID(ideaID uuid.UUID, config *config.Config, c *gin.Context) (*dto.GetIdeaResponse, error)
	ToggleFavorite(ideaID, userID uuid.UUID) (*dto.ToggleFavoriteResponse, error)
}

type ideaService struct {
	ideaRepo repository.IdeaRepository
	userRepo repository.UserRepository
}

func NewIdeaService(ideaRepo repository.IdeaRepository, userRepo repository.UserRepository) IdeaService {
	return &ideaService{ideaRepo: ideaRepo, userRepo: userRepo}
}

func (s *ideaService) GetIdeas(query *dto.GetIdeasRequest, config *config.Config, c *gin.Context) (*dto.GetIdeasResponse, error) {
	userID, statusCode := middleware.GetUserID(config.JWTSecret, s.userRepo, c)
	if query.IsFavorite {
		// не зарегистрированный пользователь не может получить список избранных
		if statusCode != 0 {
			return nil, ErrUnauthorized
		}
	}

	ideasBlock, total, err := s.ideaRepo.GetIdeas(query, userID)
	if err != nil {
		return nil, ErrInternal
	}

	ideasResponse := dto.GetIdeasResponse{Total: total, Ideas: ideasBlock}
	return &ideasResponse, nil
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

func (s *ideaService) GetIdeaByID(ideaID uuid.UUID, config *config.Config, c *gin.Context) (*dto.GetIdeaResponse, error) {
	// получаем userID, если зарегистрирован пользователей для доп информации о идее
	userID, _ := middleware.GetUserID(config.JWTSecret, s.userRepo, c)

	ideaResponse, err := s.ideaRepo.GetIdeaByIDIncr(ideaID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIdeaNotFound
		}
		return nil, ErrInternal
	}

	return ideaResponse, nil
}

func (s *ideaService) UpdateIdeaByID(ideaID, userID uuid.UUID, updateRequest *dto.UpdateIdeaRequest) (*dto.GetIdeaResponse, error) {
	// является ли пользователь владельцем данной идеи
	isUserIdeaAuthor, err := s.ideaRepo.IsUserIdeaAuthor(ideaID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isUserIdeaAuthor {
		return nil, ErrUserNotAuthor
	}

	// обновление идеи по ID
	updatedIdea, err := s.ideaRepo.UpdateIdeaByID(ideaID, userID, updateRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrIdeaConflict
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIdeaNotFound
		}
		return nil, ErrInternal
	}

	return updatedIdea, nil
}

func CheckRightsOnIdea(ideaID uuid.UUID, userID uuid.UUID, db *gorm.DB) (bool, error) {
	var isAdmin bool
	res := db.Model(&models.User{}).Where("id = ?", userID).Select("is_admin").First(&isAdmin)

	if res.Error != nil {
		return false, res.Error
	}

	if isAdmin {
		return true, nil
	}

	var authorIDStr string
	res = db.Model(&models.Idea{}).Where("id = ?", ideaID).Select("author_id").First(&authorIDStr)

	if res.Error != nil {
		return false, res.Error
	}

	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		return false, err
	}

	return authorID == userID, nil
}

func DeleteIdeaByID(ideaID uuid.UUID, db *gorm.DB) error {
	res := db.Delete(&models.Idea{}, "id = ?", ideaID)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *ideaService) ToggleFavorite(ideaID, userID uuid.UUID) (*dto.ToggleFavoriteResponse, error) {
	isFavorite, err := s.ideaRepo.ToggleFavorite(ideaID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIdeaNotFound
		}
		return nil, ErrInternal
	}

	toggleFavoriteResponse := dto.ToggleFavoriteResponse{IsFavorite: isFavorite}

	return &toggleFavoriteResponse, nil
}
