package services

import (
	"errors"
	"fmt"
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
	UpdateProjectRequest(requestID, userID uuid.UUID, status string) (*models.ProjectRequest, error)
	GetProjectRequests(projectID, userID uuid.UUID, slotID *uuid.UUID, status *string) ([]models.ProjectRequest, error)
	GetUserRequests(userID uuid.UUID) ([]models.ProjectRequest, error)
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
		Status: "pending",
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

func (s *projectRequestService) UpdateProjectRequest(requestID, userID uuid.UUID, status string) (*models.ProjectRequest, error) {
	var updatedProjectRequest *models.ProjectRequest

	err := s.projectRequestRepo.Transaction(func(tx *gorm.DB) error {
		projectRequestRepo := s.projectRequestRepo.WithTx(tx)
		projectSlotRepo := s.projectSlotRepo.WithTx(tx)
		projectRepo := s.projectRepo.WithTx(tx)

		// получение projectID через requestID
		projectID, err := projectRequestRepo.GetProjectIDByRequestID(requestID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProjectRequestNotFound
			}
			return fmt.Errorf("%w: get project by request id failed: %v", ErrInternal, err)
		}

		isUserProjectLeader, err := projectRepo.IsUserProjectLeader(projectID, userID)
		if err != nil {
			return fmt.Errorf("%w: check project leader failed: %v", ErrInternal, err)
		}

		// получение заявки
		projectRequest, err := projectRequestRepo.GetProjectRequestByID(requestID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProjectRequestNotFound
			}
			return fmt.Errorf("%w: get project request by id failed: %v", ErrInternal, err)
		}

		// заявка должна быть в статусе ожидания
		if projectRequest.Status != "pending" {
			return ErrProjectRequestNotPending
		}

		// если тип заявки "apply" то ее может принять/отклонить только лидер
		if projectRequest.Type == "apply" && !isUserProjectLeader {
			return ErrUserNotLeader
		}

		// если тип заявки "invite", то ее может принять отклонить только пользователь заявки
		if projectRequest.Type == "invite" && projectRequest.UserID != userID {
			return ErrProjectRequestDontBelongToUser
		}

		// обновление заявки
		rowsAffected, err := projectRequestRepo.UpdateProjectRequest(requestID, status)
		if err != nil {
			return fmt.Errorf("%w: update project request failed: %v", ErrInternal, err)
		}

		if rowsAffected == 0 {
			return ErrProjectRequestNotFound
		}

		// добавление пользователя в слот проекта, если статус == "accepted"
		if status == "accepted" {
			err = projectSlotRepo.PutUserIntoSlot(projectRequest.SlotID, projectRequest.UserID)
			if err != nil {
				return fmt.Errorf("%w: put user into slot failed: %v", ErrInternal, err)
			}
		}

		projectRequest.Status = status
		updatedProjectRequest = projectRequest
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedProjectRequest, nil
}

func (s *projectRequestService) GetProjectRequests(projectID, userID uuid.UUID, slotID *uuid.UUID, status *string) ([]models.ProjectRequest, error) {
	// Проверяем, является ли пользователь лидером проекта
	isLeader, err := s.projectRepo.IsUserProjectLeader(projectID, userID)
	if err != nil {
		return nil, ErrInternal
	}
	if !isLeader {
		return nil, ErrUserNotLeader
	}

	// получаем список заявок по проекту (+ фильтры по slotID и status)
	requests, err := s.projectRequestRepo.GetProjectRequests(projectID, slotID, status)
	if err != nil {
		return nil, ErrInternal
	}
	return requests, nil
}

func (s *projectRequestService) GetUserRequests(userID uuid.UUID) ([]models.ProjectRequest, error) {
	requests, err := s.projectRequestRepo.GetUserRequests(userID)
	if err != nil {
		return nil, ErrInternal
	}
	return requests, nil
}
