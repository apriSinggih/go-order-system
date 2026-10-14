package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-order-system/internal/entity"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]entity.Product, error)
	GetByID(ctx context.Context, id int64) (*entity.Product, error)
	UpdateStock(ctx context.Context, id int64, qty int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll(ctx context.Context) ([]entity.Product, error) {
	query := `SELECT id, name, price, stock, created_at, updated_at FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	query := `SELECT id, name, price, stock, created_at, updated_at FROM products WHERE id = $1`
	var p entity.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id int64, qty int) error {
	query := `UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2 AND stock >= $1`
	res, err := r.db.ExecContext(ctx, query, qty, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("stok tidak mencukupi atau produk tidak ditemukan")
	}
	return nil
}