package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/avito-shop-service/internal/models"
	"github.com/avito-shop-service/internal/repository"
	"github.com/avito-shop-service/pkg/auth"
)

type UserServiceInterface interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, username, password string) (*models.User, error)
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type UserService struct {
	repository repository.UserRepositoryInterface
}

func NewUserService(repository repository.UserRepositoryInterface) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(ctx context.Context, username, password string) (*models.User, error) {
	existingUser, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %v", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with this username already exists")
	}

	hashedPassword := HashPassword(password)

	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Coins:    1000,
	}

	err = s.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("username mustn't be empty")
	}
	return s.repository.GetUserByUsername(ctx, username)
}

func (s *UserService) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("error fetching user: %v", err)
	}
	if user == nil || !checkPasswordHash(password, user.Password) {
		return "", fmt.Errorf("invalid username or password")
	}

	signedToken, err := auth.GenerateToken(user.Username)
	if err != nil {
		return "", fmt.Errorf("error with generating token: %v", err)
	}

	return signedToken, nil
}

func (s *UserService) UpdateUserCoins(ctx context.Context, userID int64, coins int) error {
	if coins < 0 {
		return fmt.Errorf("negative amount of coins")
	}
	return s.repository.UpdateUserCoins(ctx, userID, coins)
}

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func checkPasswordHash(password, hashedPassword string) bool {
	return HashPassword(password) == hashedPassword
}
