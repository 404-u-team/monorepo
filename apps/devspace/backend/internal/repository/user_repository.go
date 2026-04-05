package repository

import (
	"errors"
	"log"
	"sort"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(payload *dto.RegisterRequest) (uuid.UUID, error)
	GetUsersTotalCount() (int64, error)
	IsUserExistByEmail(email string) (bool, error)
	IsUserExistByID(id uuid.UUID) (bool, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByNickname(login string) (*models.User, error)
	CheckUserIsAdmin(id uuid.UUID) (bool, error)
	UpdateUserByID(userID uuid.UUID, updateRequest *dto.UpdateUserRequest) error
	GetUserSkills(userID uuid.UUID) ([]dto.SkillCategoryResponse, error)
	GetUsersByParams(
		startAt, limit *uint,
		username *string,
		mainRole *uuid.UUID,
		requiredSkills *dto.UUIDSlice,
	) ([]models.User, int64, error)
}

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &userRepository{conn: conn}
}

func (r *userRepository) CreateUser(payload *dto.RegisterRequest) (uuid.UUID, error) {
	user := models.User{
		Email:        payload.Email,
		Nickname:     payload.Nickname,
		PasswordHash: payload.Password,
	}
	result := r.conn.Create(&user)
	if result.Error != nil {
		log.Println("Ошибка при создании пользователя: ", result.Error)
		return uuid.Nil, result.Error
	}

	return user.ID, nil
}

func (r *userRepository) GetUsersTotalCount() (int64, error) {
	var count int64
	err := r.conn.Model(&models.User{}).
		Count(&count).Error
	return count, err
}

func (r *userRepository) IsUserExistByEmail(email string) (bool, error) {
	var count int64
	err := r.conn.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		log.Println("Ошибка при проверке наличия пользователя: ", err)
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) IsUserExistByID(id uuid.UUID) (bool, error) {
	var count int64
	err := r.conn.Model(&models.User{}).
		Where("id = ?", id).
		Count(&count).Error
	if err != nil {
		log.Println("Ошибка при проверке наличия пользователя: ", err)
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User

	result := r.conn.First(&user, "id = ?", userID)
	if result.Error != nil {
		log.Println("Ошибка при получении пользователя по ID: ", result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	result := r.conn.First(&user, "email = ?", email)
	if result.Error != nil {
		log.Println("Ошибка при получении пользователя по email: ", result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) GetUserByNickname(nickname string) (*models.User, error) {
	var user models.User

	result := r.conn.First(&user, "nickname = ?", nickname)
	if result.Error != nil {
		log.Println("Ошибка при получении пользователя по nickname: ", result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) CheckUserIsAdmin(id uuid.UUID) (bool, error) {
	var user models.User

	result := r.conn.Select("is_admin").Where("id = ?", id).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("Пользователь не найден по id: ", id)
			return false, nil
		}
		log.Println("Ошибка при получении прав юзера по id: ", result.Error)
		return false, result.Error
	}

	return user.IsAdmin, nil
}

// обновить nickname и bio пользователя по ID. Возвращает ошибку
func (r *userRepository) UpdateUserByID(userID uuid.UUID, updateRequest *dto.UpdateUserRequest) error {
	updates := map[string]string{}

	if updateRequest.Nickname != nil {
		updates["nickname"] = *updateRequest.Nickname
	}

	if updateRequest.Bio != nil {
		updates["bio"] = *updateRequest.Bio
	}

	result := r.conn.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		log.Println("Ошибка при обновлении юзера по id: ", result.Error)
		return result.Error
	}

	return nil
}
func (r *userRepository) GetUserSkills(userID uuid.UUID) ([]dto.SkillCategoryResponse, error) {
	// подумал, что лучше один сложный запрос вместо кучи запросов для дочерних обьектов дерева

	// находим все skillID которые есть в дереве навыков
	var skillIDs []uuid.UUID
	query := `
        WITH RECURSIVE user_skill_tree AS (
            SELECT skill_id FROM "User_Skill" WHERE user_id = $1
            UNION
            SELECT sc.id FROM "Skill_Category" sc
            JOIN user_skill_tree ust ON sc.parent_id = ust.skill_id
        )
        SELECT DISTINCT skill_id FROM user_skill_tree
    `
	result := r.conn.Raw(query, userID).Scan(&skillIDs)
	if result.Error != nil {
		log.Println("Ошибка при получении рекурсивно списка скиллов пользователя: ", result.Error)
		return nil, result.Error
	}

	if len(skillIDs) == 0 {
		return []dto.SkillCategoryResponse{}, nil
	}

	// получаем скиллы на основе ID
	var allSkills []models.SkillCategory
	result = r.conn.Where("id IN ?", skillIDs).Find(&allSkills)
	if result.Error != nil {
		log.Println("Ошибка при списка скиллов по IDs: ", result.Error)
		return nil, result.Error
	}

	// преобразовываем список скиллов в дерево (на основе полей ID и parentID)
	return r.buildTreeRecursive(allSkills), nil
}

func (r *userRepository) buildTreeRecursive(skills []models.SkillCategory) []dto.SkillCategoryResponse {
	if len(skills) == 0 {
		return []dto.SkillCategoryResponse{}
	}

	// создаем map детей для каждого родителя
	childrenMap := make(map[uuid.UUID][]models.SkillCategory)
	var roots []models.SkillCategory

	for _, skill := range skills {
		if skill.ParentID == nil {
			// Это корневой навык
			roots = append(roots, skill)
		} else {
			// Это дочерний навык
			parentID := *skill.ParentID
			childrenMap[parentID] = append(childrenMap[parentID], skill)
		}
	}

	// сортируем корни для консистентного порядка
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].Name < roots[j].Name
	})

	// рекурсивно строим дерево
	result := make([]dto.SkillCategoryResponse, 0, len(roots))

	for _, root := range roots {
		tree := buildNode(root, childrenMap)
		result = append(result, tree)
	}

	return result
}

// buildNode рекурсивно строит узел дерева
func buildNode(skill models.SkillCategory, childrenMap map[uuid.UUID][]models.SkillCategory) dto.SkillCategoryResponse {
	node := dto.SkillCategoryResponse{
		ID:       skill.ID,
		Name:     skill.Name,
		ParentID: skill.ParentID,
		Children: []dto.SkillCategoryResponse{},
	}

	// добавляем детей
	if children, exists := childrenMap[skill.ID]; exists {
		// сортируем детей по имени
		sort.Slice(children, func(i, j int) bool {
			return children[i].Name < children[j].Name
		})

		for _, child := range children {
			childNode := buildNode(child, childrenMap)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}

func (r *userRepository) GetUsersByParams(
	startAt, limit *uint,
	username *string,
	mainRole *uuid.UUID,
	requiredSkills *dto.UUIDSlice,
) ([]models.User, int64, error) {

	var users []models.User

	// подгружаем навки
	query := r.conn.Session(&gorm.Session{}).
		Model(&models.User{}).
		Preload("Skills")

	// Фильтры
	if username != nil && *username != "" {
		query = query.Where("nickname ILIKE ?", "%"+(*username)+"%")
	}

	if mainRole != nil {
		query = query.Where("main_role = ?", *mainRole)
	}

	// Фильтр по навыкам
	if requiredSkills != nil && len(*requiredSkills) > 0 {
		query = query.
			Joins("JOIN user_skills ON users.id = user_skills.user_id").
			Where("user_skills.skill_id IN ?", []uuid.UUID(*requiredSkills)).
			Group("users.id").
			Having("COUNT(DISTINCT user_skills.skill_id) = ?", len(*requiredSkills))
	}

	// получение total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Пагинация
	if startAt != nil {
		query = query.Offset(int(*startAt))
	}
	if limit != nil {
		query = query.Limit(int(*limit))
	}

	// Выполняем
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
