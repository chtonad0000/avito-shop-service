package models

import "time"

type CoinTransaction struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	CounterpartUser string    `json:"to_user"`
	Amount          int       `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}
