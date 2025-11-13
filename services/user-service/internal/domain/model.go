package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user profile
type User struct {
	ID          uint64         `json:"id" gorm:"primaryKey"`
	Email       string         `json:"email" gorm:"uniqueIndex;not null"`
	Username    string         `json:"username" gorm:"uniqueIndex;not null"`
	DisplayName string         `json:"display_name"`
	Bio         string         `json:"bio" gorm:"type:text"`
	AvatarURL   string         `json:"avatar_url"`
	Region      string         `json:"region" gorm:"type:char(2)"`
	Metadata    string         `json:"metadata" gorm:"type:jsonb"` // JSON string for flexibility
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}


