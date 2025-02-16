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

type MerchRepositoryInterface interface {
	GetMerchByID(ctx context.Context, id int64) (*models.Merch, error)
	GetMerchByName(ctx context.Context, name string) (*models.Merch, error)
	GetAllMerch(ctx context.Context) ([]models.Merch, error)
	CreateMerch(ctx context.Context, merch *models.Merch) error
}

type MerchRepository struct {
	DB *pgxpool.Pool
}

func NewMerchRepository(db *pgxpool.Pool) *MerchRepository {
	return &MerchRepository{DB: db}
}

func (r *MerchRepository) GetMerchByID(ctx context.Context, id int64) (*models.Merch, error) {
	merch := &models.Merch{}
	query := `SELECT id, item_name, price FROM merch WHERE id = $1`

	err := r.DB.QueryRow(ctx, query, id).Scan(&merch.ID, &merch.ItemName, &merch.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("error fetching merch: %v", err)
		return nil, fmt.Errorf("error fetching merch: %v", err)
	}
	return merch, nil
}

func (r *MerchRepository) GetMerchByName(ctx context.Context, name string) (*models.Merch, error) {
	merch := &models.Merch{}
	query := `SELECT id, item_name, price FROM merch WHERE item_name = $1`

	err := r.DB.QueryRow(ctx, query, name).Scan(&merch.ID, &merch.ItemName, &merch.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("error fetching merch: %v", err)
		return nil, fmt.Errorf("error fetching merch: %v", err)
	}
	return merch, nil
}

func (r *MerchRepository) CreateMerch(ctx context.Context, merch *models.Merch) error {
	query := `INSERT INTO merch (item_name, price) VALUES ($1, $2)`
	_, err := r.DB.Exec(ctx, query, merch.ItemName, merch.Price)
	if err != nil {
		log.Printf("error creating merch: %v", err)
		return fmt.Errorf("error creating merch: %v", err)
	}
	return nil
}

func (r *MerchRepository) GetAllMerch(ctx context.Context) ([]models.Merch, error) {
	var merchItems []models.Merch
	query := `SELECT id, item_name, price FROM merch`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		log.Printf("error fetching merch: %v", err)
		return nil, fmt.Errorf("error fetching merch: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		merch := models.Merch{}
		err = rows.Scan(&merch.ID, &merch.ItemName, &merch.Price)
		if err != nil {
			log.Printf("error scanning row: %v", err)
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		merchItems = append(merchItems, merch)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return merchItems, nil
}
