package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/avito-shop-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUserCoins(ctx context.Context, userID int64, coins int) error
}

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password, coins FROM users WHERE username = $1`

	err := r.DB.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("error fetching user: %v", err)
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, password, coins) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(ctx, query, user.Username, user.Password, user.Coins)
	if err != nil {
		log.Printf("error creating user: %v", err)
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (r *UserRepository) UpdateUserCoins(ctx context.Context, userID int64, coins int) error {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("error starting transaction: %v", err)
		return fmt.Errorf("error starting transaction: %v", err)
	}

	_, err = tx.Exec(ctx, `UPDATE users SET coins = $1 WHERE id = $2`, coins, userID)
	if err != nil {
		log.Printf("error updating coins for user %d: %v", userID, err)
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			log.Printf("error rolling back transaction: %v", rollbackErr)
			return fmt.Errorf("error rolling back transaction: %v, original error: %v", rollbackErr, err)
		}
		return fmt.Errorf("error updating coins for user %d: %v", userID, err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("error committing transaction: %v", err)
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			log.Printf("error rolling back transaction after commit failure: %v", rollbackErr)
			return fmt.Errorf("error committing transaction: %v, error rolling back: %v", err, rollbackErr)
		}
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
