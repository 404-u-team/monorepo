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

type ProjectRequestService interface {
	CreateProjectRequestApply(payload *dto.CreateProjectRequestApplyRequest, slotID, userID, projectID uuid.UUID) (*models.ProjectRequest, error)
	CreateProjectRequestInvite(payload *dto.CreateProjectRequestInviteRequest, slotID, userID, projectID uuid.UUID) (*models.ProjectRequest, error)
}

type projectRequestService struct {
	projectRequestRepo repository.ProjectRequestRepository
	projectSlotRepo    repository.SlotRepository
	projectRepo        repository.ProjectRepository
}

func NewProjectRequestService(projectRequestRepo repository.ProjectRequestRepository, projectSlotRepo repository.SlotRepository, projectRepo repository.ProjectRepository) ProjectRequestService {
	return &projectRequestService{
		projectRequestRepo: projectRequestRepo,
		projectSlotRepo:    projectSlotRepo,
		projectRepo:        projectRepo,
	}
}

func (s *projectRequestService) CreateProjectRequestApply(payload *dto.CreateProjectRequestApplyRequest, slotID, userID, projectID uuid.UUID) (*models.ProjectRequest, error) {
	// проверка наличия проекта и слота
	slotExists, err := s.projectSlotRepo.IsSlotExists(slotID)
	if err != nil {
		return nil, ErrInternal
	}
	if !slotExists {
		return nil, ErrSlotNotFound
	}

	isSlotBelongToProject, err := s.projectSlotRepo.IsSlotBelongToProject(slotID, projectID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isSlotBelongToProject {
		return nil, ErrProjectNotFound
	}

	// проверка занят ли слот
	isSlotOpen, err := s.projectSlotRepo.IsSlotOpen(slotID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSlotNotFound
		}
		return nil, ErrInternal
	}
	if !isSlotOpen {
		return nil, ErrSlotIsClosed
	}

	// лидер проекта не может подать заявку на свой проект
	isLeader, err := s.projectRepo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if isLeader {
		return nil, ErrUserLeader
	}

	// Создание отклика на слот
	projectRequest := models.ProjectRequest{
		SlotID: slotID,
		UserID: userID,
		Type:   "apply",
		Status: "open",
	}
	if payload.CoverLetter != nil {
		projectRequest.CoverLetter = *payload.CoverLetter
	}

	if err := s.projectRequestRepo.CreateProjectRequest(&projectRequest); err != nil {
		log.Println("Ошибка при создании заявки на проект: ", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrProjectRequestConflict
		}
		return nil, ErrInternal
	}

	return &projectRequest, nil
}

func (s *projectRequestService) CreateProjectRequestInvite(payload *dto.CreateProjectRequestInviteRequest, slotID, userID, projectID uuid.UUID) (*models.ProjectRequest, error) {
	payloadUserID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Println("Ошибка при конвертировании userID в UUID")
		return nil, ErrBadRequest
	}

	// нельзя пригласить самого себя
	if userID == payloadUserID {
		return nil, ErrCantInviteYourself
	}

	// проверка наличия слота
	slotExists, err := s.projectSlotRepo.IsSlotExists(slotID)
	if err != nil {
		return nil, ErrInternal
	}
	if !slotExists {
		return nil, ErrSlotNotFound
	}

	// пользователь должен быть лидером для создания приглашения
	isLeader, err := s.projectRepo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isLeader {
		return nil, ErrUserNotLeader
	}

	// проверка, что слот принадлежит этому проекту
	isSlotBelongToProject, err := s.projectSlotRepo.IsSlotBelongToProject(slotID, projectID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isSlotBelongToProject {
		return nil, ErrProjectNotFound
	}

	// проверка занят ли слот
	isSlotOpen, err := s.projectSlotRepo.IsSlotOpen(slotID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSlotNotFound
		}
		return nil, ErrInternal
	}
	if !isSlotOpen {
		return nil, ErrSlotIsClosed
	}

	// Создание отклика на слот
	projectRequest := models.ProjectRequest{
		SlotID: slotID,
		UserID: payloadUserID,
		Type:   "invite",
		Status: "open",
	}
	if payload.CoverLetter != nil {
		projectRequest.CoverLetter = *payload.CoverLetter
	}

	if err := s.projectRequestRepo.CreateProjectRequest(&projectRequest); err != nil {
		log.Println("Ошибка при создании заявки на проект: ", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrProjectRequestConflict
		}
		return nil, ErrInternal
	}

	return &projectRequest, nil
}
