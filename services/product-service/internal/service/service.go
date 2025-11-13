package service

import (
	"errors"
	"time"

	"github.com/my-username/billion-user-app/pkg/kafkaclient"
	"github.com/my-username/billion-user-app/services/product-service/internal/domain"
	"github.com/my-username/billion-user-app/services/product-service/internal/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrUnauthorized    = errors.New("unauthorized")
)

// ProductService defines the interface for product business logic
type ProductService interface {
	CreateProduct(product *domain.Product) (*domain.Product, error)
	GetProductByID(id uint64) (*domain.Product, error)
	UpdateProduct(id uint64, updates *domain.Product, requesterID uint64) (*domain.Product, error)
	DeleteProduct(id uint64, requesterID uint64) error
	ListProducts(offset, limit int) ([]*domain.Product, error)
	SearchProducts(query string, limit int) ([]*domain.Product, error)
	GetProductsByCategory(category string, offset, limit int) ([]*domain.Product, error)
}

type productService struct {
	repo        repository.ProductRepository
	kafkaClient *kafkaclient.Client
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository, kafkaClient *kafkaclient.Client) ProductService {
	return &productService{
		repo:        repo,
		kafkaClient: kafkaClient,
	}
}

func (s *productService) CreateProduct(product *domain.Product) (*domain.Product, error) {
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	// Publish event
	event := kafkaclient.ProductCreatedEvent{
		ProductID: product.ID,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt.Format(time.RFC3339),
	}
	_ = s.kafkaClient.PublishEvent("product.created", event)

	return product, nil
}

func (s *productService) GetProductByID(id uint64) (*domain.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) UpdateProduct(id uint64, updates *domain.Product, requesterID uint64) (*domain.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check authorization (only creator can update)
	if product.CreatedBy != requesterID {
		return nil, ErrUnauthorized
	}

	// Update fields
	if updates.Name != "" {
		product.Name = updates.Name
	}
	if updates.Description != "" {
		product.Description = updates.Description
	}
	if updates.Price > 0 {
		product.Price = updates.Price
	}
	if updates.SKU != "" {
		product.SKU = updates.SKU
	}
	if updates.Stock >= 0 {
		product.Stock = updates.Stock
	}
	if updates.Category != "" {
		product.Category = updates.Category
	}
	if updates.ImageURL != "" {
		product.ImageURL = updates.ImageURL
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProduct(id uint64, requesterID uint64) error {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if product.CreatedBy != requesterID {
		return ErrUnauthorized
	}

	return s.repo.Delete(id)
}

func (s *productService) ListProducts(offset, limit int) ([]*domain.Product, error) {
	return s.repo.List(offset, limit)
}

func (s *productService) SearchProducts(query string, limit int) ([]*domain.Product, error) {
	return s.repo.Search(query, limit)
}

func (s *productService) GetProductsByCategory(category string, offset, limit int) ([]*domain.Product, error) {
	return s.repo.GetByCategory(category, offset, limit)
}
