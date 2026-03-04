package services

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"
)

type ProjectService interface {
	CreateProject(payload *dto.CreateProjectRequest, leaderID uuid.UUID) error
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) *projectService {
	return &projectService{repo: repo}
}

func (s *projectService) CreateProject(payload *dto.CreateProjectRequest, leaderID uuid.UUID) error {
	exists, err := s.repo.IsProjectExistsByTitle(payload.Title)
	if err != nil {
		return ErrInternal
	}
	if exists {
		return ErrProjectConflict
	}

	project := &models.Project{
		LeaderId:    leaderID,
		Title:       payload.Title,
		Description: payload.Description,
		Status:      "In Progress", // пока что так, может прийдется оставить
	}
	if payload.IdeaID != nil {
		project.IdeaId = *payload.IdeaID
	}
	if err := s.repo.CreateProject(project); err != nil {
		return ErrInternal
	}

	return nil
}
