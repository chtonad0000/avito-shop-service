package models

type Inventory struct {
	ID       int64 `json:"id"`
	UserID   int64 `json:"user_id"`
	ItemID   int64 `json:"item_id"`
	Quantity int   `json:"quantity"`
}
