package services

import (
	"errors"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/dto"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SlotService interface {
	GetSlots(projectID uuid.UUID) ([]models.ProjectSlot, error)
	CreateSlot(projectID, userID uuid.UUID, payload *dto.CreateSlotRequest) error
	UpdateSlotByID(slotID, projectID, userID uuid.UUID, updateRequest *dto.UpdateSlotRequest) error
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

func (s *slotService) CreateSlot(projectID, userID uuid.UUID, payload *dto.CreateSlotRequest) error {
	// является ли пользователь владельцем данного проекта
	isUserProjectLeader, err := s.projectRepo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return ErrInternal
	}
	if !isUserProjectLeader {
		return ErrUserNotLeader
	}

	// создание слота
	slot := &models.ProjectSlot{
		ProjectID:       projectID,
		SkillCategoryID: payload.SkillCategoryID,
		Title:           payload.Title,
		Status:          "open",
	}
	if payload.Description != nil {
		slot.Description = payload.Description
	}

	err = s.slotRepo.CreateSlot(projectID, slot)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrSlotConflict
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return ErrInternal
	}
	return nil
}

func (s *slotService) UpdateSlotByID(slotID, projectID, userID uuid.UUID, updateRequest *dto.UpdateSlotRequest) error {
	// проверка на пустой payload
	if updateRequest.SkillCategoryID == nil && updateRequest.Title == nil && updateRequest.Description == nil && updateRequest.Status == nil {
		return ErrEmptyPayload
	}

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
		return ErrUserNotLeader
	}

	// обновление слота по ID
	rowsAffected, err := s.slotRepo.UpdateSlotByID(slotID, projectID, updateRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrSlotConflict
		}
		return ErrInternal
	}

	if rowsAffected == 0 {
		return ErrSlotNotFound
	}

	return nil
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
		return ErrUserNotLeader
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
