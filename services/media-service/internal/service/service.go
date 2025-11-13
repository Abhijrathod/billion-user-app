package service

import (
	"errors"

	"github.com/my-username/billion-user-app/services/media-service/internal/domain"
	"github.com/my-username/billion-user-app/services/media-service/internal/repository"
)

var (
	ErrMediaNotFound = errors.New("media not found")
	ErrUnauthorized  = errors.New("unauthorized")
)

// MediaService defines the interface for media business logic
type MediaService interface {
	CreateMedia(media *domain.Media) (*domain.Media, error)
	GetMediaByID(id uint64) (*domain.Media, error)
	GetMediaByUserID(userID uint64, offset, limit int) ([]*domain.Media, error)
	DeleteMedia(id uint64, requesterID uint64) error
	GeneratePresignedURL(bucket, key string, expiresIn int) (string, error)
}

type mediaService struct {
	repo repository.MediaRepository
	// In production, add S3 client here
}

// NewMediaService creates a new media service
func NewMediaService(repo repository.MediaRepository) MediaService {
	return &mediaService{repo: repo}
}

func (s *mediaService) CreateMedia(media *domain.Media) (*domain.Media, error) {
	if err := s.repo.Create(media); err != nil {
		return nil, err
	}
	return media, nil
}

func (s *mediaService) GetMediaByID(id uint64) (*domain.Media, error) {
	return s.repo.GetByID(id)
}

func (s *mediaService) GetMediaByUserID(userID uint64, offset, limit int) ([]*domain.Media, error) {
	return s.repo.GetByUserID(userID, offset, limit)
}

func (s *mediaService) DeleteMedia(id uint64, requesterID uint64) error {
	media, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if media.UserID != requesterID {
		return ErrUnauthorized
	}

	return s.repo.Delete(id)
}

// GeneratePresignedURL generates a presigned URL for S3 upload/download
// This is a placeholder - in production, integrate with AWS SDK or MinIO
func (s *mediaService) GeneratePresignedURL(bucket, key string, expiresIn int) (string, error) {
	// TODO: Implement S3 presigned URL generation
	// For now, return a placeholder
	return "https://s3.example.com/" + bucket + "/" + key, nil
}

