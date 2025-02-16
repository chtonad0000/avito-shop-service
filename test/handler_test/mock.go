package handler

import (
	"context"
	"github.com/avito-shop-service/internal/middleware"
	"github.com/avito-shop-service/internal/models"
	"github.com/stretchr/testify/mock"
)

func setEmployeeUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, middleware.EmployeeUsernameKey, username)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, username string, password string) (*models.User, error) {
	args := m.Called(ctx, username, password)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) Authenticate(ctx context.Context, username string, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

type MockMerchService struct {
	mock.Mock
}

func (m *MockMerchService) GetMerchByName(ctx context.Context, name string) (*models.Merch, error) {
	args := m.Called(ctx, name)
	if merch, ok := args.Get(0).(*models.Merch); ok {
		return merch, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchService) GetMerchByID(ctx context.Context, id int64) (*models.Merch, error) {
	args := m.Called(ctx, id)
	if merch, ok := args.Get(0).(*models.Merch); ok {
		return merch, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchService) GetAllMerch(ctx context.Context) ([]models.Merch, error) {
	args := m.Called(ctx)
	if merch, ok := args.Get(0).([]models.Merch); ok {
		return merch, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchService) CreateMerch(ctx context.Context, merch *models.Merch) error {
	args := m.Called(ctx, merch)
	return args.Error(0)
}

type MockInventoryService struct {
	mock.Mock
}

func (m *MockInventoryService) BuyItemToInventory(ctx context.Context, userID int64, itemID int64, quantity int, price int) error {
	args := m.Called(ctx, userID, itemID, quantity, price)
	return args.Error(0)
}

func (m *MockInventoryService) GetInventoryByUserID(ctx context.Context, userID int64) ([]models.Inventory, error) {
	args := m.Called(ctx, userID)
	if inv, ok := args.Get(0).([]models.Inventory); ok {
		return inv, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInventoryService) UpdateItemQuantity(ctx context.Context, userID int, itemID int, quantity int) error {
	args := m.Called(ctx, userID, itemID, quantity)
	return args.Error(0)
}

func (m *MockInventoryService) GetItemFromInventory(ctx context.Context, userID int, itemID int) (*models.Inventory, error) {
	args := m.Called(ctx, userID, itemID)
	if inv, ok := args.Get(0).(*models.Inventory); ok {
		return inv, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInventoryService) RemoveItemFromInventory(ctx context.Context, userID int, itemID int) error {
	args := m.Called(ctx, userID, itemID)
	return args.Error(0)
}

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(ctx context.Context, transaction *models.CoinTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionService) GetTransactionsByUserId(ctx context.Context, userID int64) ([]models.CoinTransaction, error) {
	args := m.Called(ctx, userID)
	if transactions, ok := args.Get(0).([]models.CoinTransaction); ok {
		return transactions, args.Error(1)
	}
	return nil, args.Error(1)
}
