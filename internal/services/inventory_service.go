package services

import (
	"context"
	"fmt"
	"github.com/avito-shop-service/internal/models"
	"github.com/avito-shop-service/internal/repository"
)

type InventoryServiceInterface interface {
	GetInventoryByUserID(ctx context.Context, userID int64) ([]models.Inventory, error)
	BuyItemToInventory(ctx context.Context, userID int64, itemID int64, quantity int, price int) error
	UpdateItemQuantity(ctx context.Context, userID int, itemID int, quantity int) error
	GetItemFromInventory(ctx context.Context, userID int, itemID int) (*models.Inventory, error)
	RemoveItemFromInventory(ctx context.Context, userID int, itemID int) error
}

type InventoryService struct {
	repo repository.InventoryRepositoryInterface
}

func NewInventoryService(repo repository.InventoryRepositoryInterface) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) GetInventoryByUserID(ctx context.Context, userID int64) ([]models.Inventory, error) {
	if userID < 0 {
		return nil, fmt.Errorf("user id mustn't be negative")
	}
	inventoryItems, err := s.repo.GetInventoryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting inventory: %v", err)
	}
	return inventoryItems, nil
}

func (s *InventoryService) BuyItemToInventory(ctx context.Context, userID int64, itemID int64, quantity int, price int) error {
	if userID < 0 {
		return fmt.Errorf("user id mustn't be negative")
	}
	if itemID < 0 {
		return fmt.Errorf("item id mustn't be negative")
	}
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	err := s.repo.BuyItemToInventory(ctx, userID, itemID, quantity, price)
	if err != nil {
		return fmt.Errorf("error adding item to inventory: %v", err)
	}
	return nil
}

func (s *InventoryService) UpdateItemQuantity(ctx context.Context, userID int, itemID int, quantity int) error {
	if userID < 0 {
		return fmt.Errorf("user id mustn't be negative")
	}
	if itemID < 0 {
		return fmt.Errorf("item id mustn't be negative")
	}
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	err := s.repo.UpdateItemQuantity(ctx, userID, itemID, quantity)
	if err != nil {
		return fmt.Errorf("error updating item quantity: %v", err)
	}
	return nil
}

func (s *InventoryService) GetItemFromInventory(ctx context.Context, userID int, itemID int) (*models.Inventory, error) {
	if userID < 0 {
		return nil, fmt.Errorf("user id mustn't be negative")
	}
	if itemID < 0 {
		return nil, fmt.Errorf("item id mustn't be negative")
	}
	item, err := s.repo.GetItemFromInventory(ctx, userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("error fetching item from inventory: %v", err)
	}
	return item, nil
}

func (s *InventoryService) RemoveItemFromInventory(ctx context.Context, userID int, itemID int) error {
	if userID < 0 {
		return fmt.Errorf("user id mustn't be negative")
	}
	if itemID < 0 {
		return fmt.Errorf("item id mustn't be negative")
	}
	err := s.repo.RemoveItemFromInventory(ctx, userID, itemID)
	if err != nil {
		return fmt.Errorf("error removing item from inventory: %v", err)
	}
	return nil
}
