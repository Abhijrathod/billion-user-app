package domain

import (
	"time"

	"gorm.io/gorm"
)

// MediaType represents media type
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeFile  MediaType = "file"
)

// Media represents a media file
type Media struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	UserID      uint64    `json:"user_id" gorm:"not null;index"`
	FileName    string    `json:"file_name" gorm:"not null"`
	FileType    MediaType `json:"file_type" gorm:"type:varchar(20)"`
	FileSize    int64     `json:"file_size"` // in bytes
	URL         string    `json:"url" gorm:"not null"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Metadata    string    `json:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (Media) TableName() string {
	return "media"
}

