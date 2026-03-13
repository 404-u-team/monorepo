package services

import (
	"errors"
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectService interface {
	CreateProject(payload *dto.CreateProjectRequest, leaderID uuid.UUID) (*models.Project, error)
	GetProjects(query *dto.GetProjectsQuery) ([]models.Project, error)
	GetProjectByID(projectID uuid.UUID) (*models.Project, error)
	UpdateProjectByID(projectID, userID uuid.UUID, updateRequest *dto.UpdateProjectRequest) error
	DeleteProjectByID(projectID, userID uuid.UUID) error
	GetProjectRequests(projectID, userID uuid.UUID, slotID *uuid.UUID, status *string) ([]models.Request, error)
	GetUserRequests(userID uuid.UUID) ([]models.Request, error)
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) CreateProject(payload *dto.CreateProjectRequest, leaderID uuid.UUID) (*models.Project, error) {
	// проверка на наличие уже проекта с таким названием
	exists, err := s.repo.IsProjectExistsByTitle(payload.Title)
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
	if err := s.repo.CreateProject(project); err != nil {
		log.Println("Ошибка при создании проекта: ", err)
		return nil, ErrInternal
	}

	return project, nil
}

func (s *projectService) GetProjects(query *dto.GetProjectsQuery) ([]models.Project, error) {
	projects, err := s.repo.GetProjects(query)
	if err != nil {
		return nil, ErrInternal
	}

	return projects, nil
}

func (s *projectService) GetProjectByID(projectID uuid.UUID) (*models.Project, error) {
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, ErrInternal
	}

	return project, nil
}
func (s *projectService) UpdateProjectByID(projectID, userID uuid.UUID, updateRequest *dto.UpdateProjectRequest) error {
	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := s.repo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return ErrInternal
	}
	if !isUserProjectLeader {
		return ErrUserNotLeader
	}

	// обновление проекта по ID
	rowsAffected, err := s.repo.UpdateProjectbyID(projectID, updateRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrProjectConflict
		}
		return ErrInternal
	}

	if rowsAffected == 0 {
		return ErrProjectNotFound
	}

	return nil
}

func (s *projectService) DeleteProjectByID(projectID, userID uuid.UUID) error {
	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := s.repo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return ErrInternal
	}
	if !isUserProjectLeader {
		return ErrUserNotLeader
	}

	// удаление проекта по ID
	status, err := s.repo.DeleteProjectByID(projectID)
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

func (s *projectService) GetProjectRequests(projectID, userID uuid.UUID, slotID *uuid.UUID, status *string) ([]models.Request, error) {
	// Проверяем, является ли пользователь лидером проекта
	isLeader, err := s.repo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isLeader {
		return nil, ErrUserNotLeader
	}

	requests, err := s.repo.GetProjectRequests(projectID, slotID, status)
	if err != nil {
		return nil, ErrInternal
	}
	return requests, nil
}

func (s *projectService) GetUserRequests(userID uuid.UUID) ([]models.Request, error) {
	requests, err := s.repo.GetUserRequests(userID)
	if err != nil {
		return nil, ErrInternal
	}
	return requests, nil
}
