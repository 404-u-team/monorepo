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
	CreateSlot(projectID uuid.UUID, payload *dto.CreateSlotRequest) error
}

type slotService struct {
	repo repository.SlotRepository
}

func NewSlotService(repo repository.SlotRepository) SlotService {
	return &slotService{repo: repo}
}

func (s *slotService) GetSlots(projectID uuid.UUID) ([]models.ProjectSlot, error) {
	projectSlots, err := s.repo.GetSlots(projectID)
	if err != nil {
		return nil, ErrInternal
	}
	return projectSlots, nil
}

func (s *slotService) CreateSlot(projectID uuid.UUID, payload *dto.CreateSlotRequest) error {
	slot := &models.ProjectSlot{
		ProjectId:       projectID,
		SkillCategoryId: payload.SkillCategoryID,
		Title:           payload.Title,
		Status:          "open",
	}
	if payload.Description != nil {
		slot.Description = payload.Description
	}

	err := s.repo.CreateSlot(projectID, slot)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrSlotConflict
		}
		return ErrInternal
	}
	return nil
}
