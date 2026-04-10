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

type ProjectService interface {
	CreateProject(payload *dto.CreateProjectRequest, leaderID uuid.UUID) (*dto.GetProjectResponse, error)
	GetProjects(query *dto.GetProjectsQuery, config *config.Config, c *gin.Context) (*dto.GetProjectsResponse, error)
	GetProjectByID(projectID uuid.UUID, config *config.Config, c *gin.Context) (*dto.GetProjectResponse, error)
	UpdateProjectByID(projectID, userID uuid.UUID, updateRequest *dto.UpdateProjectRequest) (*dto.GetProjectResponse, error)
	DeleteProjectByID(projectID, userID uuid.UUID) error
	ToggleFavorite(projectID, userID uuid.UUID) (*dto.ToggleFavoriteResponse, error)
}

type projectService struct {
	projectRepo repository.ProjectRepository
	userRepo    repository.UserRepository
}

func NewProjectService(projectRepo repository.ProjectRepository, userRepo repository.UserRepository) ProjectService {
	return &projectService{projectRepo: projectRepo, userRepo: userRepo}
}

func (s *projectService) CreateProject(payload *dto.CreateProjectRequest, leaderID uuid.UUID) (*dto.GetProjectResponse, error) {
	// проверка на наличие уже проекта с таким названием
	exists, err := s.projectRepo.IsProjectExistsByTitle(payload.Title)
	if err != nil {
		return nil, ErrInternal
	}
	if exists {
		return nil, ErrProjectConflict
	}

	// Создание проекта
	project := &models.Project{
		LeaderID: leaderID,
		Title:    payload.Title,
		Status:   "open",
	}
	if payload.Description != nil {
		project.Description = payload.Description
	}
	if payload.IdeaID != nil {
		project.IdeaID = payload.IdeaID
	}
	if payload.Content != nil {
		project.Content = payload.Content
	}

	projectResponse, err := s.projectRepo.CreateProject(project)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrProjectConflict
		}
		return nil, ErrInternal
	}

	return projectResponse, nil
}

func (s *projectService) GetProjects(query *dto.GetProjectsQuery, config *config.Config, c *gin.Context) (*dto.GetProjectsResponse, error) {
	userID, _ := middleware.GetUserID(config.JWTSecret, s.userRepo, c)

	projectsBlock, total, err := s.projectRepo.GetProjects(query, userID)
	if err != nil {
		return nil, ErrInternal
	}

	projectsResponse := dto.GetProjectsResponse{Total: total, Projects: projectsBlock}

	return &projectsResponse, nil
}

func (s *projectService) GetProjectByID(projectID uuid.UUID, config *config.Config, c *gin.Context) (*dto.GetProjectResponse, error) {
	// получаем userID, если зарегистрирован пользователей для доп информации о идее
	userID, _ := middleware.GetUserID(config.JWTSecret, s.userRepo, c)

	projectResponse, err := s.projectRepo.GetProjectByIDIncr(projectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, ErrInternal
	}

	return projectResponse, nil
}

func (s *projectService) UpdateProjectByID(projectID, userID uuid.UUID, updateRequest *dto.UpdateProjectRequest) (*dto.GetProjectResponse, error) {
	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := s.projectRepo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isUserProjectLeader {
		return nil, ErrUserNotLeader
	}

	// обновление проекта по ID
	rowsAffected, err := s.projectRepo.UpdateProjectbyID(projectID, updateRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrProjectConflict
		}
		return nil, ErrInternal
	}

	if rowsAffected == 0 {
		return nil, ErrProjectNotFound
	}

	// получаем обновленный проект для возвращения
	projectResponse, err := s.projectRepo.GetProjectByID(projectID, userID)
	if err != nil {
		return nil, ErrInternal
	}

	return projectResponse, nil
}

func (s *projectService) DeleteProjectByID(projectID, userID uuid.UUID) error {
	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := s.projectRepo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return ErrInternal
	}
	if !isUserProjectLeader {
		return ErrUserNotLeader
	}

	// удаление проекта по ID
	status, err := s.projectRepo.DeleteProjectByID(projectID)
	if err != nil {
		return ErrInternal
	}
	if status == -1 {
		return ErrProjectHasSlots
	}
	if status == 0 {
		return ErrProjectNotFound
	}

	return nil
}

func (s *projectService) ToggleFavorite(projectID, userID uuid.UUID) (*dto.ToggleFavoriteResponse, error) {
	isFavorite, err := s.projectRepo.ToggleFavorite(projectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, ErrInternal
	}

	toggleFavoriteResponse := dto.ToggleFavoriteResponse{IsFavorite: isFavorite}

	return &toggleFavoriteResponse, nil
}
