package service

import (
	"errors"
	"time"

	"github.com/my-username/billion-user-app/pkg/kafkaclient"
	"github.com/my-username/billion-user-app/services/user-service/internal/domain"
	"github.com/my-username/billion-user-app/services/user-service/internal/repository"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUnauthorized      = errors.New("unauthorized")
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUserByID(id uint64) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	UpdateUser(id uint64, updates *domain.User, requesterID uint64) (*domain.User, error)
	DeleteUser(id uint64, requesterID uint64) error
	ListUsers(offset, limit int) ([]*domain.User, error)
	SearchUsers(query string, limit int) ([]*domain.User, error)
}

type userService struct {
	repo        repository.UserRepository
	kafkaClient *kafkaclient.Client
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, kafkaClient *kafkaclient.Client) UserService {
	return &userService{
		repo:        repo,
		kafkaClient: kafkaClient,
	}
}

func (s *userService) CreateUser(user *domain.User) (*domain.User, error) {
	// Check if user already exists
	_, err := s.repo.GetByEmail(user.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if err != repository.ErrUserNotFound {
		return nil, err
	}

	_, err = s.repo.GetByUsername(user.Username)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if err != repository.ErrUserNotFound {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Publish event
	event := kafkaclient.UserCreatedEvent{
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	_ = s.kafkaClient.PublishEvent("user.created", event)

	return user, nil
}

func (s *userService) GetUserByID(id uint64) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetUserByUsername(username string) (*domain.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *userService) UpdateUser(id uint64, updates *domain.User, requesterID uint64) (*domain.User, error) {
	// Check authorization (users can only update their own profile)
	if id != requesterID {
		return nil, ErrUnauthorized
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if updates.DisplayName != "" {
		user.DisplayName = updates.DisplayName
	}
	if updates.Bio != "" {
		user.Bio = updates.Bio
	}
	if updates.AvatarURL != "" {
		user.AvatarURL = updates.AvatarURL
	}
	if updates.Region != "" {
		user.Region = updates.Region
	}
	if updates.Metadata != "" {
		user.Metadata = updates.Metadata
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	// Publish event
	event := kafkaclient.UserUpdatedEvent{
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
	_ = s.kafkaClient.PublishEvent("user.updated", event)

	return user, nil
}

func (s *userService) DeleteUser(id uint64, requesterID uint64) error {
	if id != requesterID {
		return ErrUnauthorized
	}
	return s.repo.Delete(id)
}

func (s *userService) ListUsers(offset, limit int) ([]*domain.User, error) {
	return s.repo.List(offset, limit)
}

func (s *userService) SearchUsers(query string, limit int) ([]*domain.User, error) {
	return s.repo.Search(query, limit)
}


