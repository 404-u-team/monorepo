package repository

import (
	"fmt"
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	IsProjectExistsByTitle(title string) (bool, error)
	CreateProject(project *models.Project) error
	GetProjects(query *dto.GetProjectsQuery) ([]models.Project, error)
	GetProjectByID(projectID uuid.UUID) (*models.Project, error)
	GetProjectByTitle(title string) (*models.Project, error)
	UpdateProjectbyID(projectID uuid.UUID, updateRequest *dto.UpdateProjectRequest) (int, error)
	DeleteProjectByID(projectID uuid.UUID) (int, error)
	IsUserProjectLeader(projectID, userID uuid.UUID) (bool, error)
	GetProjectRequests(projectID uuid.UUID, slotID *uuid.UUID, status *string) ([]dto.SafeRequest, error)
	GetUserRequests(userID uuid.UUID) ([]dto.SafeRequest, error)
}

type projectRepository struct {
	conn *gorm.DB
}

func NewProjectRepository(conn *gorm.DB) ProjectRepository {
	return &projectRepository{conn: conn}
}

func (r *projectRepository) IsProjectExistsByTitle(title string) (bool, error) {
	var count int64
	result := r.conn.Model(&models.Project{}).Where("title = ?", title).Count(&count)
	if result.Error != nil {
		log.Println("Ошибка при проверке наличия проекта по названию: ", result.Error)
		return false, result.Error
	}

	return count > 0, nil
}

func (r *projectRepository) CreateProject(project *models.Project) error {
	var count int64
	if err := r.conn.Model(&models.User{}).
		Where("id = ?", project.LeaderID).
		Count(&count).Error; err != nil {
		log.Println("Ошибка при проверке наличия проекта: ", err)
		return err
	}
	if count == 0 {
		return fmt.Errorf("не найден пользователей с таким ID, невозможно создать проекта")
	}

	result := r.conn.Create(project)
	if result.Error != nil {
		log.Println("Ошибка при создании проекта: ", result.Error)
		return result.Error
	}

	return nil
}

func (r *projectRepository) GetProjects(query *dto.GetProjectsQuery) ([]models.Project, error) {
	var projects []models.Project
	result := r.conn.Model(&models.Project{})
	if query.Status != nil {
		if *query.Status == "open" || *query.Status == "closed" {
			result = result.Where("status = ?", *query.Status)
		}
	}
	if query.LeaderID != nil {
		result = result.Where("leader_id = ?", *query.LeaderID)
	}
	if query.Search != nil && *query.Search != "" {
		result = result.Where("title ILIKE ?", "%"+*query.Search+"%")
	}
	if query.StartAt != nil && *query.StartAt > 0 {
		result = result.Offset(*query.StartAt)
	}

	if query.Limit != nil {
		if *query.Limit > 50 {
			*query.Limit = 50
		}
		if *query.Limit >= 0 {
			result = result.Limit(*query.Limit)
		}
	} else {
		result = result.Limit(50)
	}

	result = result.Find(&projects)
	if result.Error != nil {
		log.Println("Ошибка при получения списка проектов: ", result.Error)
		return nil, result.Error
	}

	return projects, nil
}

func (r *projectRepository) GetProjectByID(projectID uuid.UUID) (*models.Project, error) {
	var project models.Project
	result := r.conn.First(&project, "id = ?", projectID)
	if result.Error != nil {
		log.Println("Ошибка при получении проекта по ID: ", result.Error)
		return nil, result.Error
	}

	return &project, nil
}

func (r *projectRepository) GetProjectByTitle(title string) (*models.Project, error) {
	var project models.Project
	result := r.conn.First(&project, "title = ?", title)
	if result.Error != nil {
		log.Println("Ошибка при получении проекта по title: ", result.Error)
		return nil, result.Error
	}

	return &project, nil
}

// обновить данные проекта, возвращает кол-во измененных строк и ошибку
func (r *projectRepository) UpdateProjectbyID(projectID uuid.UUID, updateRequest *dto.UpdateProjectRequest) (int, error) {
	updates := map[string]string{}

	if updateRequest.Title != nil {
		updates["title"] = *updateRequest.Title
	}

	if updateRequest.Description != nil {
		updates["description"] = *updateRequest.Description
	}

	if updateRequest.Status != nil {
		updates["status"] = *updateRequest.Status
	}

	result := r.conn.Model(&models.Project{}).Where("id = ?", projectID).Updates(updates)
	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, nil
	}

	return int(result.RowsAffected), nil
}

// возвращает статус (-1: есть слоты связанные с проектов, 0: нет проекта, 1: удален) и ошибку
func (r *projectRepository) DeleteProjectByID(projectID uuid.UUID) (int, error) {
	// проверка на наличие занятых слотов в проекте
	var count int64
	result := r.conn.Model(&models.ProjectSlot{}).Where("project_id = ?", projectID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	// есть слоты привязанные к проекту
	if count != 0 {
		return -1, nil
	}

	result = r.conn.Delete(&models.Project{}, "id = ?", projectID)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (r *projectRepository) IsUserProjectLeader(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	result := r.conn.Model(&models.Project{}).Where("id = ?", projectID).Where("leader_id = ?", userID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count == 1, nil
}

func (r *projectRepository) GetProjectRequests(projectID uuid.UUID, slotID *uuid.UUID, status *string) ([]dto.SafeRequest, error) {
	var requests []models.ProjectRequest

	query := r.conn.
		Model(&models.ProjectRequest{}). // важно!
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "nickname", "avatar_url", "main_role", "bio", "is_admin")
		}).
		Preload("Slot").
		Preload("Slot.Skill").
		Preload("Slot.Project").
		Joins("Slot").
		Where("Slot.project_id = ?", projectID)

	if slotID != nil {
		query = query.Where("request.slot_id = ?", *slotID)
	}
	if status != nil && *status != "" {
		query = query.Where("request.status = ?", *status)
	}

	if err := query.Find(&requests).Error; err != nil {
		return nil, err
	}

	return dto.FromRequests(requests), nil
}

func (r *projectRepository) GetUserRequests(userID uuid.UUID) ([]dto.SafeRequest, error) {
	var requests []models.ProjectRequest

	err := r.conn.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "nickname", "avatar_url", "main_role", "bio", "is_admin")
		}).
		Preload("Slot").
		Preload("Slot.Project", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "description", "leader_id", "status", "created_at")
		}).
		Preload("Slot.Project.Leader", func(db *gorm.DB) *gorm.DB {
			// Лидер проекта тоже User, ограничиваем поля
			return db.Select("id", "email", "nickname", "avatar_url", "main_role", "bio", "is_admin")
		}).
		Preload("Slot.Skill").
		Where("user_id = ?", userID).
		Find(&requests).Error

	if err != nil {
		return nil, err
	}

	return dto.FromRequests(requests), nil
}
