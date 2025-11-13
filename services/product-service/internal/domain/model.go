package domain

import (
	"time"

	"gorm.io/gorm"
)

// Product represents a product
type Product struct {
	ID          uint64         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"not null"`
	SKU         string         `json:"sku" gorm:"uniqueIndex"`
	Stock       int            `json:"stock" gorm:"default:0"`
	Category    string         `json:"category"`
	ImageURL    string         `json:"image_url"`
	CreatedBy   uint64         `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (Product) TableName() string {
	return "products"
}
