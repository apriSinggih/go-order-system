package service

import (
	"context"
	"testing"
	"time"

	"go-order-system/internal/entity"
)

func TestNotificationService_SendOrderPaidNotification(t *testing.T) {
	repo := newMockUserRepo()
	ctx := context.Background()

	// Seed user
	user := &entity.User{
		Email: "customer@example.com",
	}
	_ = repo.Create(ctx, user)

	notifService := NewNotificationService(repo)

	// 1. Success with existing user
	event := entity.OrderPaidEvent{
		OrderID:     101,
		UserID:      user.ID,
		TotalAmount: 250000,
		Timestamp:   time.Now(),
	}
	if err := notifService.SendOrderPaidNotification(ctx, event); err != nil {
		t.Fatalf("unexpected error sending notification: %v", err)
	}

	// 2. Fallback when user not found
	unknownEvent := entity.OrderPaidEvent{
		OrderID:     102,
		UserID:      999,
		TotalAmount: 150000,
		Timestamp:   time.Now(),
	}
	if err := notifService.SendOrderPaidNotification(ctx, unknownEvent); err != nil {
		t.Fatalf("unexpected error sending notification for unknown user: %v", err)
	}

	// 3. Fallback when userRepo is nil
	nilRepoService := NewNotificationService(nil)
	if err := nilRepoService.SendOrderPaidNotification(ctx, event); err != nil {
		t.Fatalf("unexpected error with nil userRepo: %v", err)
	}
}

