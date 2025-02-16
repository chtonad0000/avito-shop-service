package services

import (
	"context"
	"fmt"
	"github.com/avito-shop-service/internal/models"
	"github.com/avito-shop-service/internal/repository"
)

type TransactionServiceInterface interface {
	GetTransactionsByUserId(ctx context.Context, userID int64) ([]models.CoinTransaction, error)
	CreateTransaction(ctx context.Context, transaction *models.CoinTransaction) error
}

type TransactionService struct {
	repository repository.TransactionRepositoryInterface
}

func NewTransactionService(repo repository.TransactionRepositoryInterface) *TransactionService {
	return &TransactionService{repository: repo}
}

func (s *TransactionService) GetTransactionsByUserId(ctx context.Context, userID int64) ([]models.CoinTransaction, error) {
	if userID < 0 {
		return nil, fmt.Errorf("id mustn't be negative")
	}
	return s.repository.GetTransactionsByUserID(ctx, userID)
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *models.CoinTransaction) error {
	if transaction.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if transaction.UserID < 0 {
		return fmt.Errorf("id mustn't be negative")
	}

	if transaction.TransactionType != "received" && transaction.TransactionType != "send" {
		return fmt.Errorf("wrong transaction type, must be send or received")
	}

	return s.repository.CreateTransaction(ctx, transaction)
}
