//go:build unit
// +build unit

package services

import (
	"context"
	"fmt"
	"github.com/avito-shop-service/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/avito-shop-service/internal/models"
)

func TestGetMerchByID(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	merchID := int64(1)
	expectedMerch := &models.Merch{ID: merchID, ItemName: "T-Shirt", Price: 500}

	mockRepo.On("GetMerchByID", mock.Anything, merchID).Return(expectedMerch, nil)

	merch, err := service.GetMerchByID(context.Background(), merchID)
	assert.Nil(t, err)
	assert.Equal(t, expectedMerch, merch)

	mockRepo.AssertExpectations(t)
}

func TestGetMerchByID_NegativeID(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	merch, err := service.GetMerchByID(context.Background(), -1)
	assert.Nil(t, merch)
	assert.EqualError(t, err, "id mustn't be negative")
}

func TestGetMerchByName(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	merchName := "T-Shirt"
	expectedMerch := &models.Merch{ID: 1, ItemName: merchName, Price: 500}

	mockRepo.On("GetMerchByName", mock.Anything, merchName).Return(expectedMerch, nil)

	merch, err := service.GetMerchByName(context.Background(), merchName)
	assert.Nil(t, err)
	assert.Equal(t, expectedMerch, merch)

	mockRepo.AssertExpectations(t)
}

func TestGetMerchByName_NotFound(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	mockRepo.On("GetMerchByName", mock.Anything, "Unknown").Return((*models.Merch)(nil), nil)

	merch, err := service.GetMerchByName(context.Background(), "Unknown")
	assert.Nil(t, merch)
	assert.EqualError(t, err, "merch not found")
}

func TestGetAllMerch(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	expectedMerch := []models.Merch{
		{ID: 1, ItemName: "T-Shirt", Price: 500},
		{ID: 2, ItemName: "Hoodie", Price: 1200},
	}

	mockRepo.On("GetAllMerch", mock.Anything).Return(expectedMerch, nil)

	merchList, err := service.GetAllMerch(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, expectedMerch, merchList)

	mockRepo.AssertExpectations(t)
}

func TestCreateMerch(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	newMerch := &models.Merch{ID: 3, ItemName: "Cap", Price: 300}

	mockRepo.On("CreateMerch", mock.Anything, newMerch).Return(nil)

	err := service.CreateMerch(context.Background(), newMerch)
	assert.Nil(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCreateMerch_NegativePrice(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	newMerch := &models.Merch{ID: 4, ItemName: "Bag", Price: -500}

	err := service.CreateMerch(context.Background(), newMerch)
	assert.EqualError(t, err, "price mustn't be negative")
}
func TestGetMerchByID_NotFound(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	merchID := int64(999)

	mockRepo.On("GetMerchByID", mock.Anything, merchID).Return((*models.Merch)(nil), fmt.Errorf("merch not found"))

	merch, err := service.GetMerchByID(context.Background(), merchID)
	assert.Nil(t, merch)
	assert.EqualError(t, err, "merch not found")

	mockRepo.AssertExpectations(t)
}

func TestGetMerchByName_RepositoryError(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	merchName := "Hoodie"

	mockRepo.On("GetMerchByName", mock.Anything, merchName).Return((*models.Merch)(nil), fmt.Errorf("database error"))

	merch, err := service.GetMerchByName(context.Background(), merchName)
	assert.Nil(t, merch)
	assert.EqualError(t, err, "database error")

	mockRepo.AssertExpectations(t)
}

func TestGetAllMerch_EmptyList(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	mockRepo.On("GetAllMerch", mock.Anything).Return([]models.Merch{}, nil)

	merchList, err := service.GetAllMerch(context.Background())
	assert.Nil(t, err)
	assert.Empty(t, merchList)

	mockRepo.AssertExpectations(t)
}

func TestCreateMerch_RepositoryError(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	service := services.NewMerchService(mockRepo)

	newMerch := &models.Merch{ID: 5, ItemName: "Socks", Price: 150}

	mockRepo.On("CreateMerch", mock.Anything, newMerch).Return(fmt.Errorf("insert error"))

	err := service.CreateMerch(context.Background(), newMerch)
	assert.EqualError(t, err, "insert error")

	mockRepo.AssertExpectations(t)
}
