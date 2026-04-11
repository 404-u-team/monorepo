package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	IsProjectExistsByTitle(title string) (bool, error)
	CreateProject(project *models.Project) (*dto.GetProjectResponse, error)
	GetProjects(query *dto.GetProjectsQuery, userID uuid.UUID) ([]dto.ProjectBlock, int64, error)
	GetProjectByID(projectID, userID uuid.UUID) (*dto.GetProjectResponse, error)
	GetProjectByIDIncr(projectID, userID uuid.UUID) (*dto.GetProjectResponse, error)
	GetProjectByTitle(title string) (*models.Project, error)
	UpdateProjectbyID(projectID uuid.UUID, updateRequest *dto.UpdateProjectRequest) (int, error)
	DeleteProjectByID(projectID uuid.UUID) (int, error)
	ToggleFavorite(projectID, userID uuid.UUID) (bool, error)
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

func (r *projectRepository) CreateProject(project *models.Project) (*dto.GetProjectResponse, error) {
	var count int64
	if err := r.conn.Model(&models.User{}).
		Where("id = ?", project.LeaderID).
		Count(&count).Error; err != nil {
		log.Println("Ошибка при проверке наличия проекта: ", err)
		return nil, err
	}
	if count == 0 {
		err := fmt.Errorf("не найден пользователь с таким ID, невозможно создать проект")
		log.Println("Ошибка при создании проекта (не найден пользователь с таким ID): ", err)
		return nil, err
	}

	result := r.conn.Create(project)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, gorm.ErrDuplicatedKey
		}
		log.Println("Ошибка при создании проекта: ", result.Error)
		return nil, result.Error
	}

	projectResponse, err := r.GetProjectByID(project.ID, project.LeaderID)
	if err != nil {
		log.Println("Ошибка при создании проекта (получение projectResponse): ", err)
		return nil, err
	}

	return projectResponse, nil
}

func (r *projectRepository) GetProjects(query *dto.GetProjectsQuery, userID uuid.UUID) ([]dto.ProjectBlock, int64, error) {
	result := r.conn.Table("Project").
		Select(`id, leader_id, leader_id = ? AS is_leader, COALESCE("User_Favorite_Project".project_id IS NOT NULL, false) AS is_favorite, 
		title, description, views_count, favorites_count, status, idea_id, created_at, updated_at`, userID).
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

func (r *projectRepository) GetProjectByID(projectID, userID uuid.UUID) (*dto.GetProjectResponse, error) {
	var projectResponse dto.GetProjectResponse
	result := r.conn.Table("Project").Select(`id, leader_id, leader_id = ? AS is_leader, COALESCE("User_Favorite_Project".project_id IS NOT NULL, false) AS is_favorite, 
		title, description, content, COALESCE(views_count, 0) AS views_count, COALESCE(favorites_count, 0) AS favorites_count, status, idea_id, created_at, updated_at`, userID).
		Joins(`LEFT JOIN "User_Favorite_Project" ON "User_Favorite_Project".project_id = "Project".id AND "User_Favorite_Project".user_id = ?`, userID).
		Where("id = ?", projectID).
		First(&projectResponse)

	if result.Error != nil {
		log.Println("Ошибка при получении проекта по ID: ", result.Error)
		return nil, result.Error
	}

	slots, err := r.loadProjectSlots(projectID)
	if err != nil {
		return nil, err
	}
	projectResponse.Slots = slots

	return &projectResponse, nil
}

// такая же версия как и GetProjectByID, но с инкрементацией счетчика views
func (r *projectRepository) GetProjectByIDIncr(projectID, userID uuid.UUID) (*dto.GetProjectResponse, error) {
	var projectResponse dto.GetProjectResponse
	result := r.conn.Raw(`
		UPDATE "Project" p
		SET views_count = views_count + 1
		FROM (
			SELECT
				p2.id,
				p2.leader_id = ? AS is_leader,
				COALESCE(uf.project_id IS NOT NULL, false) AS is_favorite
			FROM "Project" p2
			LEFT JOIN "User_Favorite_Project" uf
				ON uf.project_id = p2.id AND uf.user_id = ?
			WHERE p2.id = ?
		) AS sub
		WHERE p.id = sub.id
		RETURNING
			p.id,
			p.leader_id,
			sub.is_leader,
			sub.is_favorite,
			p.idea_id,
			p.title,
			p.description,
			p.content,
			COALESCE(p.views_count, 0) AS views_count,
			COALESCE(p.favorites_count, 0) AS favorites_count,
			p.status,
			p.created_at,
			p.updated_at
	`, userID, userID, projectID).
		First(&projectResponse)

	if result.Error != nil {
		log.Println("Ошибка при получении проекта по ID: ", result.Error)
		return nil, result.Error
	}

	slots, err := r.loadProjectSlots(projectID)
	if err != nil {
		return nil, err
	}
	projectResponse.Slots = slots

	return &projectResponse, nil
}

func (r *projectRepository) loadProjectSlots(projectID uuid.UUID) ([]dto.GetSlotResponse, error) {
	slots, err := r.getProjectSlots(projectID)
	if err != nil {
		return nil, err
	}

	return r.mapSlotsToResponse(slots)
}

func (r *projectRepository) getProjectSlots(projectID uuid.UUID) ([]models.ProjectSlot, error) {
	var slots []models.ProjectSlot
	result := r.conn.Model(&models.ProjectSlot{}).Where("project_id = ?", projectID).Find(&slots)
	if result.Error != nil {
		log.Println("Ошибка при получении слотов проекта: ", result.Error)
		return nil, result.Error
	}

	return slots, nil
}

// вайб код
func (r *projectRepository) mapSlotsToResponse(slots []models.ProjectSlot) ([]dto.GetSlotResponse, error) {
	if len(slots) == 0 {
		return []dto.GetSlotResponse{}, nil
	}

	allSkillIDsSet := make(map[uuid.UUID]struct{})
	for _, slot := range slots {
		for _, id := range slot.PrimarySkillsID {
			allSkillIDsSet[id] = struct{}{}
		}
		for _, id := range slot.SecondarySkillsID {
			allSkillIDsSet[id] = struct{}{}
		}
	}

	allSkillIDs := make([]uuid.UUID, 0, len(allSkillIDsSet))
	for id := range allSkillIDsSet {
		allSkillIDs = append(allSkillIDs, id)
	}

	skillByID := make(map[uuid.UUID]dto.SkillCategoryResponse)
	if len(allSkillIDs) > 0 {
		var skills []models.SkillCategory
		result := r.conn.Model(&models.SkillCategory{}).
			Select("id", "parent_id", "name", "icon", "color").
			Where("id IN ?", allSkillIDs).
			Find(&skills)
		if result.Error != nil {
			log.Println("Ошибка при получении навыков слотов проекта: ", result.Error)
			return nil, result.Error
		}

		for _, skill := range skills {
			skillByID[skill.ID] = dto.SkillCategoryResponse{
				ID:       skill.ID,
				ParentID: skill.ParentID,
				Name:     skill.Name,
				Icon:     skill.Icon,
				Color:    skill.Color,
				Children: []dto.SkillCategoryResponse{},
			}
		}
	}

	result := make([]dto.GetSlotResponse, 0, len(slots))
	for _, slot := range slots {
		primarySkills := make([]dto.SkillCategoryResponse, 0, len(slot.PrimarySkillsID))
		for _, id := range slot.PrimarySkillsID {
			if skill, ok := skillByID[id]; ok {
				primarySkills = append(primarySkills, skill)
			}
		}

		secondarySkills := make([]dto.SkillCategoryResponse, 0, len(slot.SecondarySkillsID))
		for _, id := range slot.SecondarySkillsID {
			if skill, ok := skillByID[id]; ok {
				secondarySkills = append(secondarySkills, skill)
			}
		}

		result = append(result, dto.GetSlotResponse{
			ID:              slot.ID,
			PrimarySkills:   primarySkills,
			SecondarySkills: secondarySkills,
			Title:           slot.Title,
			Description:     slot.Description,
			Status:          slot.Status,
			UserID:          slot.UserID,
			CreatedAt:       slot.CreatedAt,
		})
	}

	return result, nil
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

func (r *projectRepository) ToggleFavorite(projectID, userID uuid.UUID) (bool, error) {
	// создаем транзакцию, у которой есть два исхода
	// 		1. Если есть строчка с project_id = true – удаляем строчку
	// 		2. Если не нашли ничего, то создаем новую строчку
	var result string // d - deleted, i - inserted, n - not found
	err := r.conn.Raw(`
		WITH
		project_exists AS (SELECT EXISTS(SELECT 1 FROM "Project" WHERE id = ?) AS exists),
		deleted AS (
			DELETE FROM "User_Favorite_Project"
			WHERE user_id = ? AND project_id = ? AND (SELECT exists FROM project_exists)
			RETURNING project_id
		),
		decremented AS (
			UPDATE "Project"
			SET favorites_count = favorites_count - 1
			WHERE id IN (SELECT project_id FROM deleted)
			RETURNING 'd' AS action
		),
		inserted AS (
			INSERT INTO "User_Favorite_Project" (user_id, project_id)
			SELECT ?, ?
			WHERE (SELECT exists FROM project_exists)
			AND NOT EXISTS (SELECT 1 FROM deleted)
			RETURNING project_id
		),
		incremented AS (
			UPDATE "Project"
			SET favorites_count = favorites_count + 1
			WHERE id IN (SELECT project_id FROM inserted)
			RETURNING 'i' AS action
		)
		SELECT action FROM decremented
		UNION ALL
		SELECT action FROM incremented
		UNION ALL
		SELECT 'n' WHERE NOT (SELECT exists FROM project_exists);
	`, projectID, userID, projectID, userID, projectID).First(&result).Error

	if err != nil {
		log.Println("Произошла ошибка при toggle favorite project: ", err)
		return false, err
	}

	// не существует проект
	if result == "n" {
		return false, gorm.ErrRecordNotFound
	}

	var isFavorite bool
	if result == "i" {
		isFavorite = true
	}

	return isFavorite, nil
}
