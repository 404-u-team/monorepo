package repository

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRequestRepository interface {
	CreateProjectRequest(projectRequest *models.ProjectRequest) error
	UpdateProjectRequest(requestID uuid.UUID, status string) (int, error)
	GetProjectIDByRequestID(requestID uuid.UUID) (uuid.UUID, error)
	GetProjectRequestByID(requestID uuid.UUID) (*models.ProjectRequest, error)
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

func (r *projectRequestRepository) UpdateProjectRequest(requestID uuid.UUID, status string) (int, error) {
	updates := map[string]string{"status": status}

	result := r.conn.Model(&models.ProjectRequest{}).Where("id = ?", requestID).Updates(updates)
	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, nil
	}

	return int(result.RowsAffected), nil
}

func (r *projectRequestRepository) GetProjectIDByRequestID(requestID uuid.UUID) (uuid.UUID, error) {
	var slotID uuid.UUID
	result := r.conn.Model(&models.ProjectRequest{}).Select("slot_id").Where("id = ?", requestID).First(&slotID)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	var projectID uuid.UUID
	result = r.conn.Model(&models.ProjectSlot{}).Select("project_id").Where("id = ?", slotID).First(&projectID)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	return projectID, nil
}

func (r *projectRequestRepository) GetProjectRequestByID(requestID uuid.UUID) (*models.ProjectRequest, error) {
	var projectRequest models.ProjectRequest
	result := r.conn.First(&projectRequest, "id = ?", requestID)
	if result.Error != nil {
		log.Println("Ошибка при получении заявки по ID: ", result.Error)
		return nil, result.Error
	}

	return &projectRequest, nil
}
