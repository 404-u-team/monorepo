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
	GetSlots(projectID uuid.UUID) ([]models.ProjectSlot, error)
	CreateSlot(projectID, userID uuid.UUID, payload *dto.CreateSlotRequest) (*models.ProjectSlot, error)
	UpdateSlotByID(slotID, projectID, userID uuid.UUID, updateRequest *dto.UpdateSlotRequest) (*models.ProjectSlot, error)
	DeleteSlotByID(slotID, projectID, userID uuid.UUID) error
}

type slotService struct {
	slotRepo    repository.SlotRepository
	projectRepo repository.ProjectRepository
}

func NewSlotService(slotRepo repository.SlotRepository, projectRepo repository.ProjectRepository) SlotService {
	return &slotService{slotRepo: slotRepo, projectRepo: projectRepo}
}

func (s *slotService) GetSlots(projectID uuid.UUID) ([]models.ProjectSlot, error) {
	projectSlots, err := s.slotRepo.GetSlots(projectID)
	if err != nil {
		return nil, ErrInternal
	}
	return projectSlots, nil
}

func (s *slotService) CreateSlot(projectID, userID uuid.UUID, payload *dto.CreateSlotRequest) (*models.ProjectSlot, error) {
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
		ProjectID:         projectID,
		PrimarySkillsID:   payload.PrimarySkillsID,
		SecondarySkillsID: payload.SecondarySkillsID,
		Title:             payload.Title,
		Status:            "open",
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
	return slot, nil
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

	secondarySkillsMap := make(map[uuid.UUID]struct{}, len(payload.SecondarySkillsID))
	for _, secondarySkillID := range payload.SecondarySkillsID {
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

	// Валидация secondary skills только если они указаны
	if len(payload.SecondarySkillsID) > 0 {
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

func (s *slotService) UpdateSlotByID(slotID, projectID, userID uuid.UUID, updateRequest *dto.UpdateSlotRequest) (*models.ProjectSlot, error) {
	// проверка на пустой payload
	if len(updateRequest.PrimarySkillsID) == 0 && len(updateRequest.SecondarySkillsID) == 0 && updateRequest.Title == nil && updateRequest.Description == nil && updateRequest.Status == nil {
		return nil, ErrEmptyPayload
	}

	// валидация skills только если обновляются primary skills
	if len(updateRequest.PrimarySkillsID) > 0 {
		isValidSkills, err := s.validateCreateSlotRequest(&dto.CreateSlotRequest{
			PrimarySkillsID:   updateRequest.PrimarySkillsID,
			SecondarySkillsID: updateRequest.SecondarySkillsID,
		})
		if err != nil {
			return nil, ErrInternal
		}
		if !isValidSkills {
			return nil, ErrInvalidSlotSkills
		}
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

	return slot, nil
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
