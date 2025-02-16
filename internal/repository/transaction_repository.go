package repository

import (
	"context"
	"fmt"
	"github.com/avito-shop-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
)

type TransactionRepositoryInterface interface {
	GetTransactionsByUserID(ctx context.Context, userID int64) ([]models.CoinTransaction, error)
	CreateTransaction(ctx context.Context, transaction *models.CoinTransaction) error
}

type TransactionRepository struct {
	DB *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) GetTransactionsByUserID(ctx context.Context, userID int64) ([]models.CoinTransaction, error) {
	var transactions []models.CoinTransaction
	query := `SELECT id, user_id, counterpart_username, amount, transaction_type, created_at 
              FROM coin_transactions WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching transactions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		transaction := models.CoinTransaction{}
		err = rows.Scan(&transaction.ID, &transaction.UserID, &transaction.CounterpartUser, &transaction.Amount, &transaction.TransactionType, &transaction.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *models.CoinTransaction) error {
	var anotherUserId int64
	var currentUsername string
	err := r.DB.QueryRow(ctx, `SELECT id FROM users WHERE username= $1`, transaction.CounterpartUser).Scan(&anotherUserId)
	if err != nil {
		log.Printf("failed to scan another user id: %v", err)
		return fmt.Errorf("failed to scan another user id: %v", err)
	}

	err = r.DB.QueryRow(ctx, `SELECT username FROM users WHERE id= $1`, transaction.UserID).Scan(&currentUsername)
	if err != nil {
		log.Printf("failed to scan current user username id: %v", err)
		return fmt.Errorf("failed to scan current user username id: %v", err)
	}

	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("failed to start transaction: %v", err)
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if strings.Contains(err.Error(), "tx is closed") {
			return
		}
		log.Printf("error rolling back transaction: %v", err)
	}()

	updateCoinsQuery := `UPDATE users SET coins = coins + $1 WHERE id = $2`
	_, err = tx.Exec(ctx, updateCoinsQuery, -transaction.Amount, transaction.UserID)
	if err != nil {
		log.Printf("error deducting coins: %v", err)
		return fmt.Errorf("error deducting coins: %v", err)
	}
	_, err = tx.Exec(ctx, updateCoinsQuery, transaction.Amount, anotherUserId)
	if err != nil {
		log.Printf("error adding coins: %v", err)
		return fmt.Errorf("error adding coins: %v", err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO coin_transactions (user_id, counterpart_username, amount, transaction_type) 
                           VALUES ($1, $2, $3, 'sent')`, transaction.UserID, transaction.CounterpartUser, transaction.Amount)
	if err != nil {
		log.Printf("failed to log sender transaction: %v", err)
		return fmt.Errorf("failed to log sender transaction: %v", err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO coin_transactions (user_id, counterpart_username, amount, transaction_type) 
                           VALUES ($1, $2, $3, 'received')`, anotherUserId, currentUsername, transaction.Amount)
	if err != nil {
		log.Printf("failed to log receiver transaction: %v", err)
		return fmt.Errorf("failed to log receiver transaction: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("error committing transaction: %v", err)
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
