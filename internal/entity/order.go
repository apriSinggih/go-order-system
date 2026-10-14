package entity

import "time"

type Order struct {
	ID          int64       `json:"id" db:"id"`
	UserID      int64       `json:"user_id" db:"user_id"`
	TotalAmount float64     `json:"total_amount" db:"total_amount"`
	Status      string      `json:"status" db:"status"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
	Items       []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID        int64   `json:"id" db:"id"`
	OrderID   int64   `json:"order_id" db:"order_id"`
	ProductID int64   `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	Price     float64 `json:"price" db:"price"`
}