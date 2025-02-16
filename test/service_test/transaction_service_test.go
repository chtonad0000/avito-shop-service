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
	"time"

	"github.com/avito-shop-service/internal/models"
)

func TestGetTransactionsByUserID(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	userID := int64(1)
	expectedTransactions := []models.CoinTransaction{
		{ID: 1, UserID: userID, CounterpartUser: "user2", Amount: 100, TransactionType: "send", CreatedAt: time.Now()},
		{ID: 2, UserID: userID, CounterpartUser: "user3", Amount: 200, TransactionType: "received", CreatedAt: time.Now()},
	}

	mockRepo.On("GetTransactionsByUserID", mock.Anything, userID).Return(expectedTransactions, nil)

	transactions, err := service.GetTransactionsByUserId(context.Background(), userID)
	assert.Nil(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	mockRepo.AssertExpectations(t)
}

func TestGetTransactionsByUserID_NegativeID(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	userID := int64(-1)

	transactions, err := service.GetTransactionsByUserId(context.Background(), userID)
	assert.Nil(t, transactions)
	assert.EqualError(t, err, "id mustn't be negative")

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_Success(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	transaction := &models.CoinTransaction{
		ID: 1, UserID: 1, CounterpartUser: "user2", Amount: 50, TransactionType: "send", CreatedAt: time.Now(),
	}

	mockRepo.On("CreateTransaction", mock.Anything, transaction).Return(nil)

	err := service.CreateTransaction(context.Background(), transaction)
	assert.Nil(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_NegativeUserID(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	transaction := &models.CoinTransaction{
		ID: 1, UserID: -1, CounterpartUser: "user2", Amount: 50, TransactionType: "send", CreatedAt: time.Now(),
	}

	err := service.CreateTransaction(context.Background(), transaction)
	assert.EqualError(t, err, "id mustn't be negative")

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_InvalidAmount(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	transaction := &models.CoinTransaction{
		ID: 1, UserID: 1, CounterpartUser: "user2", Amount: 0, TransactionType: "send", CreatedAt: time.Now(),
	}

	err := service.CreateTransaction(context.Background(), transaction)
	assert.EqualError(t, err, "amount must be positive")

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_InvalidTransactionType(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	transaction := &models.CoinTransaction{
		ID: 1, UserID: 1, CounterpartUser: "user2", Amount: 50, TransactionType: "transfer", CreatedAt: time.Now(),
	}

	err := service.CreateTransaction(context.Background(), transaction)
	assert.EqualError(t, err, "wrong transaction type, must be send or received")

	mockRepo.AssertExpectations(t)
}

func TestGetTransactionsByUserID_NoTransactions(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	userID := int64(1)
	mockRepo.On("GetTransactionsByUserID", mock.Anything, userID).Return([]models.CoinTransaction{}, nil)

	transactions, err := service.GetTransactionsByUserId(context.Background(), userID)
	assert.Nil(t, err)
	assert.Empty(t, transactions)

	mockRepo.AssertExpectations(t)
}

func TestGetTransactionsByUserID_ErrorFromRepo(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	userID := int64(1)
	mockRepo.On("GetTransactionsByUserID", mock.Anything, userID).Return(([]models.CoinTransaction)(nil), fmt.Errorf("database error"))

	transactions, err := service.GetTransactionsByUserId(context.Background(), userID)
	assert.Nil(t, transactions)
	assert.EqualError(t, err, "database error")

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_ZeroAmount(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	transaction := &models.CoinTransaction{
		ID: 1, UserID: 1, CounterpartUser: "user2", Amount: 0, TransactionType: "send", CreatedAt: time.Now(),
	}

	err := service.CreateTransaction(context.Background(), transaction)
	assert.EqualError(t, err, "amount must be positive")

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_TooLargeAmount(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	transaction := &models.CoinTransaction{
		ID: 1, UserID: 1, CounterpartUser: "user2", Amount: 1_000_000_000, TransactionType: "send", CreatedAt: time.Now(),
	}

	mockRepo.On("CreateTransaction", mock.Anything, transaction).Return(fmt.Errorf("amount exceeds limit"))

	err := service.CreateTransaction(context.Background(), transaction)
	assert.EqualError(t, err, "amount exceeds limit")

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_ErrorFromRepo(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	transaction := &models.CoinTransaction{
		ID: 1, UserID: 1, CounterpartUser: "user2", Amount: 50, TransactionType: "send", CreatedAt: time.Now(),
	}

	mockRepo.On("CreateTransaction", mock.Anything, transaction).Return(fmt.Errorf("database error"))

	err := service.CreateTransaction(context.Background(), transaction)
	assert.EqualError(t, err, "database error")

	mockRepo.AssertExpectations(t)
}
