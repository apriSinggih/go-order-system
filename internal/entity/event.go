package entity

import "time"

type OrderPaidEvent struct {
	OrderID     int64     `json:"order_id"`
	UserID      int64     `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Timestamp   time.Time `json:"timestamp"`
}