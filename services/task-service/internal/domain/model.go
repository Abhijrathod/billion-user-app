package domain

import (
	"time"

	"gorm.io/gorm"
)

// TaskStatus represents task status
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// Task represents a task
type Task struct {
	ID          uint64         `json:"id" gorm:"primaryKey"`
	UserID      uint64         `json:"user_id" gorm:"not null;index"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Status      TaskStatus     `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Priority    int            `json:"priority" gorm:"default:0"` // 0=low, 1=medium, 2=high
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (Task) TableName() string {
	return "tasks"
}
