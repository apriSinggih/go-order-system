package service

import (
	"context"
	"fmt"
	"testing"

	"go-order-system/internal/entity"
)

type mockOrderRepo struct {
	orders map[int64]*entity.Order
}

func (m *mockOrderRepo) Create(ctx context.Context, order *entity.Order, items []entity.OrderItem) error {
	order.ID = int64(len(m.orders) + 1)
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepo) GetByID(ctx context.Context, id int64) (*entity.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, fmt.Errorf("order not found")
	}
	return o, nil
}

func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	o, ok := m.orders[id]
	if !ok {
		return fmt.Errorf("order not found")
	}
	o.Status = status
	return nil
}

func TestOrderService_GetOrderOwnership(t *testing.T) {
	repo := &mockOrderRepo{
		orders: map[int64]*entity.Order{
			1: {ID: 1, UserID: 100, TotalAmount: 50000, Status: "PAID"},
			2: {ID: 2, UserID: 200, TotalAmount: 75000, Status: "PENDING"},
		},
	}

	svc := &orderService{
		orderRepo: repo,
	}

	ctx := context.Background()

	// 1. Owner can access their own order
	order, err := svc.GetOrder(ctx, 100, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.ID != 1 || order.UserID != 100 {
		t.Errorf("order mismatch: %+v", order)
	}

	// 2. Different user CANNOT access someone else's order (Ownership isolation)
	_, err = svc.GetOrder(ctx, 999, 1)
	if err == nil {
		t.Fatalf("expected error when accessing another user's order, got nil")
	}

	// 3. Non-existent order
	_, err = svc.GetOrder(ctx, 100, 999)
	if err == nil {
		t.Fatalf("expected error for non-existent order, got nil")
	}
}

