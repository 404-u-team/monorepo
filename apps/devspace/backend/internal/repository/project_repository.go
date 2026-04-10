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
	GetProjects(query *dto.GetProjectsQuery, userID uuid.UUID) ([]dto.ProjectBlock, int64, error)
	GetProjectByID(projectID uuid.UUID) (*models.Project, error)
	GetProjectByTitle(title string) (*models.Project, error)
	UpdateProjectbyID(projectID uuid.UUID, updateRequest *dto.UpdateProjectRequest) (int, error)
	DeleteProjectByID(projectID uuid.UUID) (int, error)
	IsUserProjectLeader(projectID, userID uuid.UUID) (bool, error)
	WithTx(tx *gorm.DB) ProjectRepository
}

type projectRepository struct {
	conn *gorm.DB
}

func NewProjectRepository(conn *gorm.DB) ProjectRepository {
	return &projectRepository{conn: conn}
}

func (r *projectRepository) WithTx(tx *gorm.DB) ProjectRepository {
	return &projectRepository{conn: tx}
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
		err := fmt.Errorf("не найден пользователь с таким ID, невозможно создать проект")
		log.Println("Ошибка при создании проекта (не найден пользователь с таким ID): ", err)
		return err
	}

	result := r.conn.Create(project)
	if result.Error != nil {
		log.Println("Ошибка при создании проекта: ", result.Error)
		return result.Error
	}

	return nil
}

func (r *projectRepository) GetProjects(query *dto.GetProjectsQuery, userID uuid.UUID) ([]dto.ProjectBlock, int64, error) {
	result := r.conn.Table("Project").
		Select(`id, leader_id, is_leader = ? AS is_leader, COALESCE("User_Favorite_Project".project_id IS NOT NULL, false) AS is_favorite, 
		title, description, views_count, favorites_count, status, created_at, updated_at`, userID).
		Joins(`LEFT JOIN "User_Favorite_Project" ON "User_Favorite_Project".project_id = "Project".id AND "User_Favorite_Project".user_id = ?`, userID)

	if query.Status != nil {
		if *query.Status == "open" || *query.Status == "closed" {
			result = result.Where("status = ?", *query.Status)
		}
	}
	if query.LeaderID != nil {
		result = result.Where("leader_id = ?", query.LeaderID.UUID())
	}
	if query.Search != nil && *query.Search != "" {
		result = result.Where("title ILIKE ?", "%"+*query.Search+"%")
	}

	if query.IdeaID != nil {
		result = result.Where("idea_id = ?", query.IdeaID.UUID())
	}

	if query.OpenSlots != nil && *query.OpenSlots {
		subQuery := r.conn.Model(&models.ProjectSlot{}).
			Select("1").
			Where(`project_id = "Project".id AND status='open'`)
		result = result.Where("EXISTS (?)", subQuery)
	}

	if query.SlotsSkills != nil {
		if len(*query.SlotsSkills) == 0 {
			return []dto.ProjectBlock{}, 0, nil
		}

		// находим через subquery список id проектов у которых есть
		// все навыки из SlotsSkills (внутри primary и secondary skills)
		skills := *query.SlotsSkills

		var projectIDs []uuid.UUID

		// Convert UUIDs to PostgreSQL ARRAY format: '{uuid1,uuid2,...}'::uuid[]
		skillArray := "ARRAY["
		for i, skill := range skills {
			if i > 0 {
				skillArray += ","
			}
			skillArray += "'" + skill.String() + "'::uuid"
		}
		skillArray += "]"

		var subquery *gorm.DB
		if query.OpenSlots != nil && *query.OpenSlots {
			subquery = r.conn.Raw(`
				SELECT project_id
				FROM (
					SELECT project_id, unnest(primary_skills_id || secondary_skills_id) AS skill
					FROM "Project_Slot"
					WHERE status = 'open'
				) AS all_skills
				WHERE skill = ANY(`+skillArray+`)
				GROUP BY project_id
				HAVING COUNT(DISTINCT skill) = ?
			`, len(skills)).Scan(&projectIDs)
		} else {
			subquery = r.conn.Raw(`
				SELECT project_id
				FROM (
					SELECT project_id, unnest(primary_skills_id || secondary_skills_id) AS skill
					FROM "Project_Slot"
				) AS all_skills
				WHERE skill = ANY(`+skillArray+`)
				GROUP BY project_id
				HAVING COUNT(DISTINCT skill) = ?
			`, len(skills)).Scan(&projectIDs)
		}
		if subquery.Error != nil {
			log.Println("Ошибка при фильтрации проектов по навыкам слотов: ", subquery.Error)
			return nil, 0, subquery.Error
		}

		if len(projectIDs) > 0 {
			result = result.Where("id IN ?", projectIDs)
		} else {
			return []dto.ProjectBlock{}, 0, nil
		}
	}

	if query.MinPeople != nil || query.MaxPeople != nil {
		subQuery := r.conn.Model(&models.ProjectSlot{}).
			Select("COUNT(user_id)").
			Where(`project_id = "Project".id`)

		if query.MinPeople != nil {
			result = result.Where("(?) + 1 >= ?", subQuery, *query.MinPeople)
		}
		if query.MaxPeople != nil {
			result = result.Where("(?) + 1 <= ?", subQuery, *query.MaxPeople)
		}
	}

	var total int64
	if err := result.Count(&total).Error; err != nil {
		log.Println("Ошибка при подсчете количества проектов: ", err)
		return nil, 0, err
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

	var projects []dto.ProjectBlock
	result = result.Find(&projects)
	if result.Error != nil {
		log.Println("Ошибка при получения списка проектов: ", result.Error)
		return nil, 0, result.Error
	}

	return projects, total, nil
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
	updates := map[string]interface{}{}

	if updateRequest.Title != nil {
		updates["title"] = *updateRequest.Title
	}

	if updateRequest.Description != nil {
		updates["description"] = *updateRequest.Description
	}

	if updateRequest.Status != nil {
		updates["status"] = *updateRequest.Status
	}

	if updateRequest.Content != nil {
		updates["content"] = *updateRequest.Content
	}

	result := r.conn.Model(&models.Project{}).Where("id = ?", projectID).Updates(updates)
	if result.Error != nil {
		log.Println("Ошибка при обновлении проекта: ", result.Error)
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, nil
	}

	return int(result.RowsAffected), nil
}

// при удалении проекта каскадно удаляются и слоты и запросы на слот
func (r *projectRepository) DeleteProjectByID(projectID uuid.UUID) (int, error) {
	result := r.conn.Delete(&models.Project{}, "id = ?", projectID)
	if result.Error != nil {
		log.Println("Ошибка при удалении проекта: ", result.Error)
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (r *projectRepository) IsUserProjectLeader(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	result := r.conn.Model(&models.Project{}).Where("id = ?", projectID).Where("leader_id = ?", userID).Count(&count)
	if result.Error != nil {
		log.Println("Ошибка при проверке является ли пользователь лидером проекта: ", result.Error)
		return false, result.Error
	}
	return count == 1, nil
}
