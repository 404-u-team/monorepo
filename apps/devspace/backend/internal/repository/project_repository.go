package repository

import (
	"fmt"
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	IsProjectExistsByTitle(title string) (bool, error)
	CreateProject(project *models.Project) error
	GetProjects(query *dto.GetProjectsQuery) ([]models.Project, error)
}

type projectRepository struct {
	conn *gorm.DB
}

func NewProjectRepository(conn *gorm.DB) ProjectRepository {
	return &projectRepository{conn: conn}
}

func (r *projectRepository) IsProjectExistsByTitle(title string) (bool, error) {
	var exists bool
	err := r.conn.Model(&models.Project{}).
		Select("COUNT(*) = 1").
		Where("title = ?", title).
		Find(&exists).Error
	if err != nil {
		log.Println("Ошибка при проверке наличия проекта по названию: ", err)
		return false, err
	}

	return exists, err
}

func (r *projectRepository) CreateProject(project *models.Project) error {
	var cnt int64
	if err := r.conn.Model(&models.User{}).
		Where("id = ?", project.LeaderID).
		Count(&cnt).Error; err != nil {
		log.Println("Ошибка при проверке наличия проекта: ", err)
		return err
	}
	if cnt == 0 {
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
	if query.Limit != nil && *query.Limit >= 0 {
		result = result.Limit(*query.Limit)
	}

	result = result.Find(&projects);
	if result.Error != nil {
		log.Println("Ошибка при получения списка проектов: ", result.Error)
		return nil, result.Error
	}

	return projects, nil
}
