package models

type Merch struct {
	ID       int64  `json:"id"`
	ItemName string `json:"item_name"`
	Price    int    `json:"price"`
}
