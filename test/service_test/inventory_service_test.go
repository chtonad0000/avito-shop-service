//go:build unit
// +build unit

package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/avito-shop-service/internal/models"
	"github.com/avito-shop-service/internal/services"
)

func TestGetInventoryByUserID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	userID := int64(1)
	expectedInventory := []models.Inventory{
		{UserID: userID, ItemID: 1, Quantity: 2},
		{UserID: userID, ItemID: 2, Quantity: 5},
	}

	mockRepo.On("GetInventoryByUserID", mock.Anything, userID).Return(expectedInventory, nil)

	inventory, err := service.GetInventoryByUserID(context.Background(), userID)
	assert.Nil(t, err)
	assert.Equal(t, expectedInventory, inventory)

	mockRepo.AssertExpectations(t)
}

func TestBuyItemToInventory(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	userID := int64(1)
	itemID := int64(2)
	quantity := 3
	price := 100

	mockRepo.On("BuyItemToInventory", mock.Anything, userID, itemID, quantity, price).Return(nil)

	err := service.BuyItemToInventory(context.Background(), userID, itemID, quantity, price)
	assert.Nil(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateItemQuantity(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	userID := 1
	itemID := 2
	quantity := 10

	mockRepo.On("UpdateItemQuantity", mock.Anything, userID, itemID, quantity).Return(nil)

	err := service.UpdateItemQuantity(context.Background(), userID, itemID, quantity)
	assert.Nil(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetItemFromInventory(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	userID := 1
	itemID := 2
	expectedItem := &models.Inventory{UserID: int64(userID), ItemID: int64(itemID), Quantity: 5}

	mockRepo.On("GetItemFromInventory", mock.Anything, userID, itemID).Return(expectedItem, nil)

	item, err := service.GetItemFromInventory(context.Background(), userID, itemID)
	assert.Nil(t, err)
	assert.Equal(t, expectedItem, item)

	mockRepo.AssertExpectations(t)
}

func TestRemoveItemFromInventory(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	userID := 1
	itemID := 2

	mockRepo.On("RemoveItemFromInventory", mock.Anything, userID, itemID).Return(nil)

	err := service.RemoveItemFromInventory(context.Background(), userID, itemID)
	assert.Nil(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetInventoryByUserID_ErrorNegativeUserID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	_, err := service.GetInventoryByUserID(context.Background(), -1)
	assert.Error(t, err)
	assert.Equal(t, "user id mustn't be negative", err.Error())
}

func TestBuyItemToInventory_ErrorNegativeUserID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	err := service.BuyItemToInventory(context.Background(), -1, 1, 1, 100)
	assert.Error(t, err)
	assert.Equal(t, "user id mustn't be negative", err.Error())
}

func TestUpdateItemQuantity_ErrorNegativeUserID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	err := service.UpdateItemQuantity(context.Background(), -1, 1, 1)
	assert.Error(t, err)
	assert.Equal(t, "user id mustn't be negative", err.Error())
}

func TestGetItemFromInventory_ErrorNegativeUserID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	_, err := service.GetItemFromInventory(context.Background(), -1, 1)
	assert.Error(t, err)
	assert.Equal(t, "user id mustn't be negative", err.Error())
}

func TestRemoveItemFromInventory_ErrorNegativeUserID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	err := service.RemoveItemFromInventory(context.Background(), -1, 1)
	assert.Error(t, err)
	assert.Equal(t, "user id mustn't be negative", err.Error())
}

func TestBuyItemToInventory_ErrorNegativeItemID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	err := service.BuyItemToInventory(context.Background(), 1, -1, 1, 100)
	assert.Error(t, err)
	assert.Equal(t, "item id mustn't be negative", err.Error())
}

func TestUpdateItemQuantity_ErrorNegativeItemID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	err := service.UpdateItemQuantity(context.Background(), 1, -1, 1)
	assert.Error(t, err)
	assert.Equal(t, "item id mustn't be negative", err.Error())
}

func TestGetItemFromInventory_ErrorNegativeItemID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	_, err := service.GetItemFromInventory(context.Background(), 1, -1)
	assert.Error(t, err)
	assert.Equal(t, "item id mustn't be negative", err.Error())
}

func TestRemoveItemFromInventory_ErrorNegativeItemID(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	err := service.RemoveItemFromInventory(context.Background(), 1, -1)
	assert.Error(t, err)
	assert.Equal(t, "item id mustn't be negative", err.Error())
}

func TestBuyItemToInventory_ErrorNonPositiveQuantity(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	err := service.BuyItemToInventory(context.Background(), 1, 1, 0, 100)
	assert.Error(t, err)
	assert.Equal(t, "quantity must be positive", err.Error())
}

func TestUpdateItemQuantity_ErrorNonPositiveQuantity(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	err := service.UpdateItemQuantity(context.Background(), 1, 1, 0)
	assert.Error(t, err)
	assert.Equal(t, "quantity must be positive", err.Error())
}

func TestGetInventoryByUserID_DBError(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	mockRepo.On("GetInventoryByUserID", mock.Anything, int64(1)).Return([]models.Inventory{}, errors.New("DB error"))

	_, err := service.GetInventoryByUserID(context.Background(), 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting inventory: DB error")

	mockRepo.AssertExpectations(t)
}

func TestBuyItemToInventory_DBError(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	mockRepo.On("BuyItemToInventory", mock.Anything, int64(1), int64(1), 1, 100).Return(errors.New("DB error"))

	err := service.BuyItemToInventory(context.Background(), 1, 1, 1, 100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error adding item to inventory: DB error")

	mockRepo.AssertExpectations(t)
}

func TestUpdateItemQuantity_DBError(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	mockRepo.On("UpdateItemQuantity", mock.Anything, 1, 1, 1).Return(errors.New("DB error"))

	err := service.UpdateItemQuantity(context.Background(), 1, 1, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error updating item quantity: DB error")

	mockRepo.AssertExpectations(t)
}

func TestGetItemFromInventory_DBError(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	mockRepo.On("GetItemFromInventory", mock.Anything, 1, 1).Return((*models.Inventory)(nil), errors.New("DB error"))

	_, err := service.GetItemFromInventory(context.Background(), 1, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error fetching item from inventory: DB error")

	mockRepo.AssertExpectations(t)
}

func TestRemoveItemFromInventory_DBError(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)

	mockRepo.On("RemoveItemFromInventory", mock.Anything, 1, 1).Return(errors.New("DB error"))

	err := service.RemoveItemFromInventory(context.Background(), 1, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error removing item from inventory: DB error")

	mockRepo.AssertExpectations(t)
}
