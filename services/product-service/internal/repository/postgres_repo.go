package repository

import (
	"errors"

	"github.com/my-username/billion-user-app/services/product-service/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(product *domain.Product) error
	GetByID(id uint64) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint64) error
	List(offset, limit int) ([]*domain.Product, error)
	Search(query string, limit int) ([]*domain.Product, error)
	GetByCategory(category string, offset, limit int) ([]*domain.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id uint64) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepository) List(offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := r.db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Search(query string, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	searchPattern := "%" + query + "%"
	if err := r.db.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?",
		searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetByCategory(category string, offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := r.db.Where("category = ?", category).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
