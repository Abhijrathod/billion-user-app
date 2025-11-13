package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/my-username/billion-user-app/pkg/jwtutils"
	"github.com/my-username/billion-user-app/pkg/kafkaclient"
	"github.com/my-username/billion-user-app/services/auth-service/internal/domain"
	"github.com/my-username/billion-user-app/services/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user account is inactive")
)

// AuthService defines the interface for auth business logic
type AuthService interface {
	Register(email, username, password string) (*domain.User, error)
	Login(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(refreshToken string) error
	ValidateToken(token string) (*jwtutils.Claims, error)
	GetUserByID(id uint64) (*domain.User, error)
}

type authService struct {
	repo        repository.AuthRepository
	jwtManager  *jwtutils.JWTManager
	kafkaClient *kafkaclient.Client
}

// NewAuthService creates a new auth service
func NewAuthService(repo repository.AuthRepository, jwtManager *jwtutils.JWTManager, kafkaClient *kafkaclient.Client) AuthService {
	return &authService{
		repo:        repo,
		jwtManager:  jwtManager,
		kafkaClient: kafkaClient,
	}
}

func (s *authService) Register(email, username, password string) (*domain.User, error) {
	// Check if user already exists by email
	_, err := s.repo.GetUserByEmail(email)
	if err == nil {
		return nil, repository.ErrUserAlreadyExists
	}
	if err != repository.ErrUserNotFound {
		return nil, err
	}

	// Check username
	_, err = s.repo.GetUserByUsername(username)
	if err == nil {
		return nil, repository.ErrUserAlreadyExists
	}
	if err != repository.ErrUserNotFound {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
		IsActive: true,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Publish "user.created" event
	if s.kafkaClient != nil {
		event := kafkaclient.UserCreatedEvent{
			UserID:    user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		}
		_ = s.kafkaClient.PublishEvent("user.created", event)
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", "", ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return "", "", err
	}

	// Create secure refresh token
	refreshToken := generateSecureToken(64)
	expiry := time.Now().Add(7 * 24 * time.Hour)

	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiry,
	}

	if err := s.repo.SaveRefreshToken(rt); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	rt, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	user, err := s.repo.GetUserByID(rt.UserID)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", "", ErrUserInactive
	}

	// Generate new access token
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return "", "", err
	}

	// Rotate refresh token
	newRefresh := generateSecureToken(64)
	rt.Token = newRefresh
	rt.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	if err := s.repo.SaveRefreshToken(rt); err != nil {
		return "", "", err
	}

	// Delete old token
	_ = s.repo.DeleteRefreshToken(refreshToken)

	return accessToken, newRefresh, nil
}

func (s *authService) Logout(refreshToken string) error {
	return s.repo.DeleteRefreshToken(refreshToken)
}

func (s *authService) ValidateToken(token string) (*jwtutils.Claims, error) {
	return s.jwtManager.ValidateToken(token)
}

func (s *authService) GetUserByID(id uint64) (*domain.User, error) {
	return s.repo.GetUserByID(id)
}

// generateSecureToken creates a cryptographically secure random string
func generateSecureToken(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
