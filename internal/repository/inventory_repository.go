package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/avito-shop-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepositoryInterface interface {
	GetInventoryByUserID(ctx context.Context, userID int64) ([]models.Inventory, error)
	BuyItemToInventory(ctx context.Context, userID int64, itemID int64, quantity int, price int) error
	UpdateItemQuantity(ctx context.Context, userID int, itemID int, quantity int) error
	GetItemFromInventory(ctx context.Context, userID int, itemID int) (*models.Inventory, error)
	RemoveItemFromInventory(ctx context.Context, userID int, itemID int) error
}

type InventoryRepository struct {
	DB *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{DB: db}
}

func (r *InventoryRepository) GetInventoryByUserID(ctx context.Context, userID int64) ([]models.Inventory, error) {
	var inventoryItems []models.Inventory
	query := `SELECT id, user_id, item_id, quantity FROM inventory WHERE user_id = $1`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		log.Printf("error fetching inventory: %v", err)
		return nil, fmt.Errorf("error fetching inventory: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		inventoryItem := models.Inventory{}
		err = rows.Scan(&inventoryItem.ID, &inventoryItem.UserID, &inventoryItem.ItemID, &inventoryItem.Quantity)
		if err != nil {
			log.Printf("error scanning row: %v", err)
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		inventoryItems = append(inventoryItems, inventoryItem)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return inventoryItems, nil
}

func (r *InventoryRepository) BuyItemToInventory(ctx context.Context, userID int64, itemID int64, quantity int, merchPrice int) error {
	var userCoins int64
	err := r.DB.QueryRow(ctx, "SELECT coins FROM users WHERE id = $1", userID).Scan(&userCoins)
	if err != nil {
		log.Printf("error fetching user balance: %v", err)
		return fmt.Errorf("error fetching user balance: %v", err)
	}
	if userCoins < int64(merchPrice*quantity) {
		log.Printf("not enough coins")
		return fmt.Errorf("not enough coins")
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

	_, err = r.DB.Exec(ctx, "UPDATE users SET coins = coins - $1 WHERE id = $2", merchPrice*quantity, userID)
	if err != nil {
		log.Printf("error updating user coins: %v", err)
		return fmt.Errorf("error updating user coins: %v", err)
	}

	query := `INSERT INTO inventory (user_id, item_id, quantity) 
              VALUES ($1, $2, $3)
              ON CONFLICT (user_id, item_id) 
              DO UPDATE SET quantity = inventory.quantity + $3`
	_, err = tx.Exec(ctx, query, userID, itemID, quantity)
	if err != nil {
		log.Printf("error adding item to inventory: %v", err)
		return fmt.Errorf("error adding item to inventory: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("error committing transaction: %v", err)
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (r *InventoryRepository) UpdateItemQuantity(ctx context.Context, userID int, itemID int, quantity int) error {
	query := `UPDATE inventory SET quantity = $3 WHERE user_id = $1 AND item_id = $2`

	_, err := r.DB.Exec(ctx, query, userID, itemID, quantity)
	if err != nil {
		log.Printf("error updating item quantity: %v", err)
		return fmt.Errorf("error updating item quantity: %v", err)
	}
	return nil
}

func (r *InventoryRepository) GetItemFromInventory(ctx context.Context, userID int, itemID int) (*models.Inventory, error) {
	query := `SELECT id, user_id, item_id, quantity FROM inventory WHERE user_id=$1 AND item_id=$2`
	model := &models.Inventory{}
	err := r.DB.QueryRow(ctx, query, userID, itemID).Scan(&model.ID, &model.UserID, &model.ItemID, &model.Quantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("error fetching item: %v", err)
		return nil, fmt.Errorf("error fetching item: %v", err)
	}
	return model, nil
}

func (r *InventoryRepository) RemoveItemFromInventory(ctx context.Context, userID int, itemID int) error {
	query := `DELETE FROM inventory WHERE user_id = $1 AND item_id = $2`

	_, err := r.DB.Exec(ctx, query, userID, itemID)
	if err != nil {
		log.Printf("error removing item from inventory: %v", err)
		return fmt.Errorf("error removing item from inventory: %v", err)
	}
	return nil
}
