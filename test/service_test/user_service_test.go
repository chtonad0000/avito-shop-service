//go:build unit
// +build unit

package services

import (
	"context"
	"errors"
	"github.com/avito-shop-service/internal/models"
	"github.com/avito-shop-service/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	username := "testuser"
	password := "password123"

	mockRepo.On("GetUserByUsername", mock.Anything, username).Return((*models.User)(nil), nil)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	createdUser, err := service.CreateUser(context.Background(), username, password)

	assert.NoError(t, err)
	assert.Equal(t, username, createdUser.Username)
	assert.Equal(t, 1000, createdUser.Coins)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	username := "existinguser"
	existingUser := &models.User{Username: username}

	mockRepo.On("GetUserByUsername", mock.Anything, username).Return(existingUser, nil)

	_, err := service.CreateUser(context.Background(), username, "password")

	assert.EqualError(t, err, "user with this username already exists")
	mockRepo.AssertExpectations(t)
}

func TestAuthenticate_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	username := "testuser"
	password := "password123"
	hashedPassword := services.HashPassword(password)
	user := &models.User{Username: username, Password: hashedPassword}

	mockRepo.On("GetUserByUsername", mock.Anything, username).Return(user, nil)

	token, err := service.Authenticate(context.Background(), username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	username := "testuser"
	password := "wrongpassword"
	hashedPassword := services.HashPassword("password123")
	user := &models.User{Username: username, Password: hashedPassword}

	mockRepo.On("GetUserByUsername", mock.Anything, username).Return(user, nil)

	_, err := service.Authenticate(context.Background(), username, password)

	assert.EqualError(t, err, "invalid username or password")
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserCoins_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	userID := int64(1)
	coins := 500

	mockRepo.On("UpdateUserCoins", mock.Anything, userID, coins).Return(nil)

	err := service.UpdateUserCoins(context.Background(), userID, coins)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserCoins_NegativeAmount(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	err := service.UpdateUserCoins(context.Background(), 1, -500)

	assert.EqualError(t, err, "negative amount of coins")
	mockRepo.AssertNotCalled(t, "UpdateUserCoins")
}

func TestCreateUser_ErrorCheckingExistingUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	username := "testuser"
	password := "password123"

	mockRepo.On("GetUserByUsername", mock.Anything, username).Return((*models.User)(nil), errors.New("database error"))

	_, err := service.CreateUser(context.Background(), username, password)

	assert.EqualError(t, err, "error checking user existence: database error")
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_ErrorCreatingUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	username := "newuser"
	password := "password123"
	mockRepo.On("GetUserByUsername", mock.Anything, username).Return((*models.User)(nil), nil)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("insert error"))

	_, err := service.CreateUser(context.Background(), username, password)

	assert.EqualError(t, err, "error creating user: insert error")
	mockRepo.AssertExpectations(t)
}

func TestGetUserByUsername_EmptyUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	_, err := service.GetUserByUsername(context.Background(), "")

	assert.EqualError(t, err, "username mustn't be empty")
	mockRepo.AssertNotCalled(t, "GetUserByUsername")
}

func TestGetUserByUsername_ErrorFetchingUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	username := "testuser"
	mockRepo.On("GetUserByUsername", mock.Anything, username).Return((*models.User)(nil), errors.New("database error"))

	_, err := service.GetUserByUsername(context.Background(), username)

	assert.EqualError(t, err, "database error")
	mockRepo.AssertExpectations(t)
}

func TestAuthenticate_ErrorFetchingUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	username := "testuser"
	password := "password123"
	mockRepo.On("GetUserByUsername", mock.Anything, username).Return((*models.User)(nil), errors.New("database error"))

	_, err := service.Authenticate(context.Background(), username, password)

	assert.EqualError(t, err, "error fetching user: database error")
	mockRepo.AssertExpectations(t)
}

func TestAuthenticate_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	username := "nonexistent"
	password := "password123"
	mockRepo.On("GetUserByUsername", mock.Anything, username).Return((*models.User)(nil), nil)

	_, err := service.Authenticate(context.Background(), username, password)

	assert.EqualError(t, err, "invalid username or password")
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserCoins_ErrorUpdatingDB(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	userID := int64(1)
	coins := 500
	mockRepo.On("UpdateUserCoins", mock.Anything, userID, coins).Return(errors.New("update failed"))

	err := service.UpdateUserCoins(context.Background(), userID, coins)

	assert.EqualError(t, err, "update failed")
	mockRepo.AssertExpectations(t)
}
