package repository

import (
	"context"
	"database/sql"

	"go-order-system/internal/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order, items []entity.OrderItem) error
	GetByID(ctx context.Context, id int64) (*entity.Order, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order, items []entity.OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryOrder := `INSERT INTO orders (user_id, total_amount, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err = tx.QueryRowContext(ctx, queryOrder, order.UserID, order.TotalAmount, order.Status).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	queryItem := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id`
	for i := range items {
		items[i].OrderID = order.ID
		err = tx.QueryRowContext(ctx, queryItem, order.ID, items[i].ProductID, items[i].Quantity, items[i].Price).Scan(&items[i].ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) GetByID(ctx context.Context, id int64) (*entity.Order, error) {
	queryOrder := `SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders WHERE id = $1`
	var o entity.Order
	err := r.db.QueryRowContext(ctx, queryOrder, id).Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}

	queryItems := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := r.db.QueryContext(ctx, queryItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}