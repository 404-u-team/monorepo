package repository

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	IsProjectExistsByTitle(title string) (bool, error)
	CreateProject(project *models.Project) error
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
	result := r.conn.Create(project)
	if result.Error != nil {
		log.Println("Ошибка при создании проекта: ", result.Error)
		return result.Error
	}

	return nil
}
