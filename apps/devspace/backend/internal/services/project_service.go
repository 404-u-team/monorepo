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
	GetProjects(query *dto.GetProjectsQuery) (*dto.GetProjectsResponse, error)
	GetProjectByID(projectID uuid.UUID) (*models.Project, error)
	UpdateProjectByID(projectID, userID uuid.UUID, updateRequest *dto.UpdateProjectRequest) (*models.Project, error)
	DeleteProjectByID(projectID, userID uuid.UUID) error
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
	if payload.Content != nil {
		project.Content = payload.Content
	}
	if err := s.repo.CreateProject(project); err != nil {
		log.Println("Ошибка при создании проекта: ", err)
		return nil, ErrInternal
	}

	return project, nil
}

func (s *projectService) GetProjects(query *dto.GetProjectsQuery) (*dto.GetProjectsResponse, error) {
	projects, total, err := s.repo.GetProjects(query)
	if err != nil {
		return nil, ErrInternal
	}

	projectsResponse := dto.GetProjectsResponse{Total: total, Projects: projects}

	return &projectsResponse, nil
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

func (s *projectService) UpdateProjectByID(projectID, userID uuid.UUID, updateRequest *dto.UpdateProjectRequest) (*models.Project, error) {
	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := s.repo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isUserProjectLeader {
		return nil, ErrUserNotLeader
	}

	// обновление проекта по ID
	rowsAffected, err := s.repo.UpdateProjectbyID(projectID, updateRequest)
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
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, ErrInternal
	}

	return project, nil
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
