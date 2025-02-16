package services

import (
	"context"
	"fmt"
	"github.com/avito-shop-service/internal/models"
	"github.com/avito-shop-service/internal/repository"
)

type MerchServiceInterface interface {
	GetMerchByID(ctx context.Context, id int64) (*models.Merch, error)
	GetMerchByName(ctx context.Context, name string) (*models.Merch, error)
	GetAllMerch(ctx context.Context) ([]models.Merch, error)
	CreateMerch(ctx context.Context, merch *models.Merch) error
}

type MerchService struct {
	repository repository.MerchRepositoryInterface
}

func NewMerchService(repo repository.MerchRepositoryInterface) *MerchService {
	return &MerchService{repository: repo}
}

func (s *MerchService) GetMerchByID(ctx context.Context, id int64) (*models.Merch, error) {
	if id < 0 {
		return nil, fmt.Errorf("id mustn't be negative")
	}
	merch, err := s.repository.GetMerchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return merch, nil
}

func (s *MerchService) GetMerchByName(ctx context.Context, name string) (*models.Merch, error) {
	merch, err := s.repository.GetMerchByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if merch == nil {
		return nil, fmt.Errorf("merch not found")
	}
	return merch, nil
}

func (s *MerchService) GetAllMerch(ctx context.Context) ([]models.Merch, error) {
	return s.repository.GetAllMerch(ctx)
}

func (s *MerchService) CreateMerch(ctx context.Context, merch *models.Merch) error {
	if merch.Price < 0 {
		return fmt.Errorf("price mustn't be negative")
	}
	return s.repository.CreateMerch(ctx, merch)
}
