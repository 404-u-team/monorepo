package services

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"
)

type ProjectService interface {
	CreateProject(payload *dto.CreateProjectRequest, leaderID uuid.UUID) (*models.Project, error)
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) *projectService {
	return &projectService{repo: repo}
}

func (s *projectService) CreateProject(payload *dto.CreateProjectRequest, leaderID uuid.UUID) (*models.Project, error) {
	exists, err := s.repo.IsProjectExistsByTitle(payload.Title)
	if err != nil {
		return nil, ErrInternal
	}
	if exists {
		return nil, ErrProjectConflict
	}

	project := &models.Project{
		LeaderID:    leaderID,
		Title:       payload.Title,
		Description: payload.Description,
		Status:      "In Progress", // пока что так, может прийдется оставить
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
