package services

import (
	"context"
	"github.com/avito-shop-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) GetInventoryByUserID(ctx context.Context, userID int64) ([]models.Inventory, error) {
	args := m.Called(ctx, userID)
	if inv, ok := args.Get(0).([]models.Inventory); ok {
		return inv, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInventoryRepository) BuyItemToInventory(ctx context.Context, userID int64, itemID int64, quantity int, price int) error {
	args := m.Called(ctx, userID, itemID, quantity, price)
	return args.Error(0)
}

func (m *MockInventoryRepository) UpdateItemQuantity(ctx context.Context, userID int, itemID int, quantity int) error {
	args := m.Called(ctx, userID, itemID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetItemFromInventory(ctx context.Context, userID int, itemID int) (*models.Inventory, error) {
	args := m.Called(ctx, userID, itemID)
	if inv, ok := args.Get(0).(*models.Inventory); ok {
		return inv, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInventoryRepository) RemoveItemFromInventory(ctx context.Context, userID int, itemID int) error {
	args := m.Called(ctx, userID, itemID)
	return args.Error(0)
}

type MockMerchRepository struct {
	mock.Mock
}

func (m *MockMerchRepository) GetMerchByID(ctx context.Context, id int64) (*models.Merch, error) {
	args := m.Called(ctx, id)
	if merch, ok := args.Get(0).(*models.Merch); ok {
		return merch, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchRepository) GetMerchByName(ctx context.Context, name string) (*models.Merch, error) {
	args := m.Called(ctx, name)
	if merch, ok := args.Get(0).(*models.Merch); ok {
		return merch, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchRepository) GetAllMerch(ctx context.Context) ([]models.Merch, error) {
	args := m.Called(ctx)
	if merch, ok := args.Get(0).([]models.Merch); ok {
		return merch, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchRepository) CreateMerch(ctx context.Context, merch *models.Merch) error {
	args := m.Called(ctx, merch)
	return args.Error(0)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetTransactionsByUserID(ctx context.Context, userID int64) ([]models.CoinTransaction, error) {
	args := m.Called(ctx, userID)
	if transactions, ok := args.Get(0).([]models.CoinTransaction); ok {
		return transactions, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, transaction *models.CoinTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserCoins(ctx context.Context, userID int64, coins int) error {
	args := m.Called(ctx, userID, coins)
	return args.Error(0)
}
