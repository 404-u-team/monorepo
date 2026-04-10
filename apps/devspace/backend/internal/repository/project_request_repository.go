package repository

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRequestRepository interface {
	CreateProjectRequest(projectRequest *models.ProjectRequest) error
	UpdateProjectRequestStatus(requestID uuid.UUID, status string) (int, error)
	GetProjectIDByRequestID(requestID uuid.UUID) (uuid.UUID, error)
	GetProjectRequestByID(requestID uuid.UUID) (*models.ProjectRequest, error)
	GetProjectRequests(projectID uuid.UUID, slotID *uuid.UUID, status *string) ([]models.ProjectRequest, error)
	GetUserRequests(userID uuid.UUID) ([]models.ProjectRequest, error)
	WithTx(tx *gorm.DB) ProjectRequestRepository
	Transaction(fn func(tx *gorm.DB) error) error
}

type projectRequestRepository struct {
	conn *gorm.DB
}

func NewProjectRequestRepository(conn *gorm.DB) ProjectRequestRepository {
	return &projectRequestRepository{conn: conn}
}

func (r *projectRequestRepository) WithTx(tx *gorm.DB) ProjectRequestRepository {
	return &projectRequestRepository{conn: tx}
}

func (r *projectRequestRepository) Transaction(fn func(tx *gorm.DB) error) error {
	err := r.conn.Transaction(fn)
	if err != nil {
		return err
	}

	return nil
}

func (r *projectRequestRepository) CreateProjectRequest(projectRequest *models.ProjectRequest) error {
	result := r.conn.Create(projectRequest)
	if result.Error != nil {
		log.Println("Ошибка при создании запроса на проект: ", result.Error)
		return result.Error
	}

	return nil
}

func (r *projectRequestRepository) UpdateProjectRequestStatus(requestID uuid.UUID, status string) (int, error) {
	updates := map[string]interface{}{"status": status}

	result := r.conn.Model(&models.ProjectRequest{}).Where("id = ?", requestID).Updates(updates)
	if result.Error != nil {
		log.Println("Ошибка при обновлении project request: ", result.Error)
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, nil
	}

	return int(result.RowsAffected), nil
}

func (r *projectRequestRepository) GetProjectIDByRequestID(requestID uuid.UUID) (uuid.UUID, error) {
	var slotIDStruct struct {
		SlotID uuid.UUID
	}
	result := r.conn.Model(&models.ProjectRequest{}).Select("slot_id").Where("id = ?", requestID).First(&slotIDStruct)
	if result.Error != nil {
		log.Println("Ошибка при получении projectID через requestID (slot_id part): ", result.Error)
		return uuid.Nil, result.Error
	}

	var projectIDStruct struct {
		ProjectID uuid.UUID
	}
	result = r.conn.Model(&models.ProjectSlot{}).Select("project_id").Where("id = ?", slotIDStruct.SlotID).First(&projectIDStruct)
	if result.Error != nil {
		log.Println("Ошибка при получении projectID через requestID (project_id part): ", result.Error)
		return uuid.Nil, result.Error
	}

	return projectIDStruct.ProjectID, nil
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

func (r *projectRequestRepository) GetProjectRequests(projectID uuid.UUID, slotID *uuid.UUID, status *string) ([]models.ProjectRequest, error) {
	var requests []models.ProjectRequest

	result := r.conn.Model(&models.ProjectRequest{}).
		Joins("JOIN \"Project_Slot\" ON \"Project_Slot\".id = \"Project_Request\".slot_id").
		Where("\"Project_Slot\".project_id = ?", projectID)

	// фильтр по slotID, если не nil
	if slotID != nil {
		result = result.Where("\"Project_Request\".slot_id = ?", *slotID)
	}

	// фильтр по status, если не nil
	if status != nil {
		result = result.Where("\"Project_Request\".status = ?", *status)
	}

	result = result.Find(&requests)
	if result.Error != nil {
		log.Println("Ошибка при получении списка заявок проекта: ", result.Error)
		return nil, result.Error
	}

	return requests, nil
}

func (r *projectRequestRepository) GetUserRequests(userID uuid.UUID) ([]models.ProjectRequest, error) {
	var requests []models.ProjectRequest

	result := r.conn.Model(&models.ProjectRequest{}).Where("user_id = ?", userID).Where("type = ?", "apply").Find(&requests)

	if result.Error != nil {
		log.Println("Ошибка при получении заявок пользователя: ", result.Error)
		return nil, result.Error
	}

	return requests, nil
}
