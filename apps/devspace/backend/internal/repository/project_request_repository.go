package repository

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

type ProjectRequestRepository interface {
	CreateProjectRequest(projectRequest *models.ProjectRequest) error
}

type projectRequestRepository struct {
	conn *gorm.DB
}

func NewProjectRequestRepository(conn *gorm.DB) ProjectRequestRepository {
	return &projectRequestRepository{conn: conn}
}

func (r *projectRequestRepository) CreateProjectRequest(projectRequest *models.ProjectRequest) error {
	result := r.conn.Create(projectRequest)
	if result.Error != nil {
		log.Println("Ошибка при создании запроса на проект: ", result.Error)
		return result.Error
	}

	return nil
}
