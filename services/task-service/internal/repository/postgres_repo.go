package repository

import (
	"errors"

	"github.com/my-username/billion-user-app/services/task-service/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

// TaskRepository defines the interface for task data operations
type TaskRepository interface {
	Create(task *domain.Task) error
	GetByID(id uint64) (*domain.Task, error)
	GetByUserID(userID uint64, offset, limit int) ([]*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id uint64) error
	GetByStatus(userID uint64, status domain.TaskStatus, offset, limit int) ([]*domain.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetByID(id uint64) (*domain.Task, error) {
	var task domain.Task
	if err := r.db.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetByUserID(userID uint64, offset, limit int) ([]*domain.Task, error) {
	var tasks []*domain.Task
	if err := r.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Update(task *domain.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Task{}, id).Error
}

func (r *taskRepository) GetByStatus(userID uint64, status domain.TaskStatus, offset, limit int) ([]*domain.Task, error) {
	var tasks []*domain.Task
	if err := r.db.Where("user_id = ? AND status = ?", userID, status).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
