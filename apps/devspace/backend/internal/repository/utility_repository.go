package repository

import (
	"sort"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Удаляет все записи в таблице с полем, значение которого val
func DeleteAll(field, tableName string, val any, db *gorm.DB) error {
	res := db.Table(tableName).Delete(nil, field+" = ?", val)
	return res.Error
}

func BuildSkillTree(skills []models.SkillCategory) []dto.SkillCategoryResponse {
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
		Icon:     skill.Icon,
		Color:    skill.Color,
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
