package services

import (
	"errors"
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SlotService interface {
	GetSlots(projectID uuid.UUID) ([]dto.GetSlotResponse, error)
	CreateSlot(projectID, userID uuid.UUID, payload *dto.CreateSlotRequest) (*dto.GetSlotResponse, error)
	UpdateSlotByID(slotID, projectID, userID uuid.UUID, updateRequest *dto.UpdateSlotRequest) (*dto.GetSlotResponse, error)
	DeleteSlotByID(slotID, projectID, userID uuid.UUID) error
}

type slotService struct {
	slotRepo    repository.SlotRepository
	projectRepo repository.ProjectRepository
}

func NewSlotService(slotRepo repository.SlotRepository, projectRepo repository.ProjectRepository) SlotService {
	return &slotService{slotRepo: slotRepo, projectRepo: projectRepo}
}

func (s *slotService) GetSlots(projectID uuid.UUID) ([]dto.GetSlotResponse, error) {
	projectSlots, err := s.slotRepo.GetSlots(projectID)
	if err != nil {
		return nil, ErrInternal
	}

	responses, err := s.mapSlotsToResponse(projectSlots)
	if err != nil {
		return nil, ErrInternal
	}

	return responses, nil
}

func (s *slotService) CreateSlot(projectID, userID uuid.UUID, payload *dto.CreateSlotRequest) (*dto.GetSlotResponse, error) {
	isValidPayload, err := s.validateCreateSlotRequest(payload)
	if err != nil {
		return nil, ErrInternal
	}
	if !isValidPayload {
		return nil, ErrInvalidSlotSkills
	}

	// есть ли такой проект
	_, err = s.projectRepo.GetProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, ErrInternal
	}

	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := s.projectRepo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isUserProjectLeader {
		return nil, ErrUserNotLeader
	}

	// создание слота
	slot := &models.ProjectSlot{
		ProjectID:       projectID,
		PrimarySkillsID: models.UUIDArray(payload.PrimarySkillsID),
		Title:           payload.Title,
		Status:          "open",
	}
	if payload.SecondarySkillsID != nil {
		slot.SecondarySkillsID = models.UUIDArray(*payload.SecondarySkillsID)
	}
	if payload.Description != nil {
		slot.Description = payload.Description
	}

	err = s.slotRepo.CreateSlot(projectID, slot)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlotConflict
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, ErrInternal
	}

	response, err := s.mapSlotToResponse(slot)
	if err != nil {
		return nil, ErrInternal
	}

	return response, nil
}

func (s *slotService) validateCreateSlotRequest(payload *dto.CreateSlotRequest) (bool, error) {
	if payload == nil {
		return false, nil
	}

	// Primary skills обязательны, secondary опциональны
	if len(payload.PrimarySkillsID) == 0 {
		return false, nil
	}

	primarySkillsMap := make(map[uuid.UUID]struct{}, len(payload.PrimarySkillsID))
	for _, primarySkillID := range payload.PrimarySkillsID {
		if _, exists := primarySkillsMap[primarySkillID]; exists {
			return false, nil
		}
		primarySkillsMap[primarySkillID] = struct{}{}
	}

	secondarySkills := make([]uuid.UUID, 0)
	if payload.SecondarySkillsID != nil {
		secondarySkills = *payload.SecondarySkillsID
	}

	secondarySkillsMap := make(map[uuid.UUID]struct{}, len(secondarySkills))
	for _, secondarySkillID := range secondarySkills {
		if _, exists := secondarySkillsMap[secondarySkillID]; exists {
			return false, nil
		}
		secondarySkillsMap[secondarySkillID] = struct{}{}
	}

	primaryIDs := make([]uuid.UUID, 0, len(primarySkillsMap))
	for id := range primarySkillsMap {
		primaryIDs = append(primaryIDs, id)
	}

	secondaryIDs := make([]uuid.UUID, 0, len(secondarySkillsMap))
	for id := range secondarySkillsMap {
		secondaryIDs = append(secondaryIDs, id)
	}

	primarySkills, err := s.slotRepo.GetSkillCategoriesByIDs(primaryIDs)
	if err != nil {
		log.Printf("Ошибка при получении primary skills: %v", err)
		return false, err
	}
	if len(primarySkills) != len(primarySkillsMap) {
		log.Printf("Primary skills не найдены. Запрашивали %d, получили %d", len(primarySkillsMap), len(primarySkills))
		return false, nil
	}

	for _, primarySkill := range primarySkills {
		if primarySkill.ParentID != nil {
			log.Printf("Primary skill %s имеет родителя, должен быть root skill", primarySkill.ID)
			return false, nil
		}
	}

	// валидация secondary skills только если они указаны
	if len(secondarySkills) > 0 {
		secondarySkills, err := s.slotRepo.GetSkillCategoriesByIDs(secondaryIDs)
		if err != nil {
			return false, err
		}
		if len(secondarySkills) != len(secondarySkillsMap) {
			return false, nil
		}

		for _, secondarySkill := range secondarySkills {
			if secondarySkill.ParentID == nil {
				return false, nil
			}
			if _, ok := primarySkillsMap[*secondarySkill.ParentID]; !ok {
				return false, nil
			}
		}
	}

	return true, nil
}

func (s *slotService) UpdateSlotByID(slotID, projectID, userID uuid.UUID, updateRequest *dto.UpdateSlotRequest) (*dto.GetSlotResponse, error) {
	// проверка на пустой payload
	if updateRequest.PrimarySkillsID == nil && updateRequest.SecondarySkillsID == nil && updateRequest.Title == nil && updateRequest.Description == nil && updateRequest.Status == nil {
		return nil, ErrEmptyPayload
	}

	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := s.projectRepo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isUserProjectLeader {
		return nil, ErrUserNotLeader
	}

	// принадлежит ли слот данному проекту
	isSlotBelongsToProject, err := s.slotRepo.IsSlotBelongToProject(slotID, projectID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isSlotBelongsToProject {
		return nil, ErrSlotNotFound
	}

	currentSlot, err := s.slotRepo.GetSlotByID(slotID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSlotNotFound
		}
		return nil, ErrInternal
	}

	// копируем скиллы из слота
	mergedPrimarySkills := make([]uuid.UUID, len(currentSlot.PrimarySkillsID))
	copy(mergedPrimarySkills, currentSlot.PrimarySkillsID)

	mergedSecondarySkills := make([]uuid.UUID, len(currentSlot.SecondarySkillsID))
	copy(mergedSecondarySkills, currentSlot.SecondarySkillsID)

	// если переданы какие-то скилы, то добавляем их
	isSkillsUpdated := false
	if updateRequest.PrimarySkillsID != nil {
		mergedPrimarySkills = *updateRequest.PrimarySkillsID
		isSkillsUpdated = true
	}
	if updateRequest.SecondarySkillsID != nil {
		mergedSecondarySkills = *updateRequest.SecondarySkillsID
		isSkillsUpdated = true
	}

	if isSkillsUpdated {
		isValidSkills, err := s.validateCreateSlotRequest(&dto.CreateSlotRequest{
			PrimarySkillsID:   mergedPrimarySkills,
			SecondarySkillsID: &mergedSecondarySkills,
		})
		if err != nil {
			return nil, ErrInternal
		}
		if !isValidSkills {
			return nil, ErrInvalidSlotSkills
		}
	}

	// обновление слота по ID
	rowsAffected, err := s.slotRepo.UpdateSlotByID(slotID, projectID, updateRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlotConflict
		}
		return nil, ErrInternal
	}

	if rowsAffected == 0 {
		return nil, ErrSlotNotFound
	}

	slot, err := s.slotRepo.GetSlotByID(slotID)
	if err != nil {
		return nil, ErrInternal
	}

	response, err := s.mapSlotToResponse(slot)
	if err != nil {
		return nil, ErrInternal
	}

	return response, nil
}

func (s *slotService) mapSlotToResponse(slot *models.ProjectSlot) (*dto.GetSlotResponse, error) {
	responses, err := s.mapSlotsToResponse([]models.ProjectSlot{*slot})
	if err != nil {
		return nil, err
	}

	if len(responses) == 0 {
		return nil, nil
	}

	return &responses[0], nil
}

func (s *slotService) mapSlotsToResponse(slots []models.ProjectSlot) ([]dto.GetSlotResponse, error) {
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
		skills, err := s.slotRepo.GetSkillCategoriesByIDs(allSkillIDs)
		if err != nil {
			return nil, err
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

func (s *slotService) DeleteSlotByID(slotID, projectID, userID uuid.UUID) error {
	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := s.projectRepo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return ErrInternal
	}
	if !isUserProjectLeader {
		return ErrUserNotLeader
	}

	// принадлежит ли слот данному проекту
	isSlotBelongsToProject, err := s.slotRepo.IsSlotBelongToProject(slotID, projectID)
	if err != nil {
		return ErrInternal
	}
	if !isSlotBelongsToProject {
		return ErrSlotNotFound
	}

	// удаление слота по ID
	rowsAffected, err := s.slotRepo.DeleteSlotByID(slotID, projectID)
	if err != nil {
		return ErrInternal
	}
	if rowsAffected == 0 {
		return ErrSlotNotFound
	}

	return nil
}
