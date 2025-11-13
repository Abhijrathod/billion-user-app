package service

import (
	"errors"
	"time"

	"github.com/my-username/billion-user-app/pkg/kafkaclient"
	"github.com/my-username/billion-user-app/services/task-service/internal/domain"
	"github.com/my-username/billion-user-app/services/task-service/internal/repository"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrUnauthorized = errors.New("unauthorized")
)

// TaskService defines the interface for task business logic
type TaskService interface {
	CreateTask(task *domain.Task) (*domain.Task, error)
	GetTaskByID(id uint64) (*domain.Task, error)
	GetTasksByUserID(userID uint64, offset, limit int) ([]*domain.Task, error)
	UpdateTask(id uint64, updates *domain.Task, requesterID uint64) (*domain.Task, error)
	DeleteTask(id uint64, requesterID uint64) error
	GetTasksByStatus(userID uint64, status domain.TaskStatus, offset, limit int) ([]*domain.Task, error)
}

type taskService struct {
	repo        repository.TaskRepository
	kafkaClient *kafkaclient.Client
}

// NewTaskService creates a new task service
func NewTaskService(repo repository.TaskRepository, kafkaClient *kafkaclient.Client) TaskService {
	return &taskService{
		repo:        repo,
		kafkaClient: kafkaClient,
	}
}

func (s *taskService) CreateTask(task *domain.Task) (*domain.Task, error) {
	if task.Status == "" {
		task.Status = domain.TaskStatusPending
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	// Publish event
	event := kafkaclient.TaskCreatedEvent{
		TaskID:    task.ID,
		UserID:    task.UserID,
		Title:     task.Title,
		Status:    string(task.Status),
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
	}
	_ = s.kafkaClient.PublishEvent("task.created", event)

	return task, nil
}

func (s *taskService) GetTaskByID(id uint64) (*domain.Task, error) {
	return s.repo.GetByID(id)
}

func (s *taskService) GetTasksByUserID(userID uint64, offset, limit int) ([]*domain.Task, error) {
	return s.repo.GetByUserID(userID, offset, limit)
}

func (s *taskService) UpdateTask(id uint64, updates *domain.Task, requesterID uint64) (*domain.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check authorization
	if task.UserID != requesterID {
		return nil, ErrUnauthorized
	}

	// Update fields
	if updates.Title != "" {
		task.Title = updates.Title
	}
	if updates.Description != "" {
		task.Description = updates.Description
	}
	if updates.Status != "" {
		task.Status = updates.Status
	}
	if updates.Priority >= 0 {
		task.Priority = updates.Priority
	}
	if updates.DueDate != nil {
		task.DueDate = updates.DueDate
	}

	if err := s.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) DeleteTask(id uint64, requesterID uint64) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if task.UserID != requesterID {
		return ErrUnauthorized
	}

	return s.repo.Delete(id)
}

func (s *taskService) GetTasksByStatus(userID uint64, status domain.TaskStatus, offset, limit int) ([]*domain.Task, error) {
	return s.repo.GetByStatus(userID, status, offset, limit)
}
