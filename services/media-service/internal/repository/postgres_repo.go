package repository

import (
	"errors"

	"github.com/my-username/billion-user-app/services/media-service/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrMediaNotFound = errors.New("media not found")
)

// MediaRepository defines the interface for media data operations
type MediaRepository interface {
	Create(media *domain.Media) error
	GetByID(id uint64) (*domain.Media, error)
	GetByUserID(userID uint64, offset, limit int) ([]*domain.Media, error)
	Delete(id uint64) error
}

type mediaRepository struct {
	db *gorm.DB
}

// NewMediaRepository creates a new media repository
func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(media *domain.Media) error {
	return r.db.Create(media).Error
}

func (r *mediaRepository) GetByID(id uint64) (*domain.Media, error) {
	var media domain.Media
	if err := r.db.First(&media, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMediaNotFound
		}
		return nil, err
	}
	return &media, nil
}

func (r *mediaRepository) GetByUserID(userID uint64, offset, limit int) ([]*domain.Media, error) {
	var media []*domain.Media
	if err := r.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&media).Error; err != nil {
		return nil, err
	}
	return media, nil
}

func (r *mediaRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Media{}, id).Error
}

