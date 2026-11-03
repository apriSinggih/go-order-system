package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-order-system/internal/entity"
	"go-order-system/internal/repository"
)

type NotificationService interface {
	SendOrderPaidNotification(ctx context.Context, event entity.OrderPaidEvent) error
}

type notificationService struct {
	userRepo repository.UserRepository
}

func NewNotificationService(userRepo repository.UserRepository) NotificationService {
	return &notificationService{
		userRepo: userRepo,
	}
}

func (s *notificationService) SendOrderPaidNotification(ctx context.Context, event entity.OrderPaidEvent) error {
	recipient := fmt.Sprintf("user_%d@system.local", event.UserID)

	if s.userRepo != nil {
		user, err := s.userRepo.GetByID(ctx, event.UserID)
		if err == nil && user != nil && user.Email != "" {
			recipient = user.Email
		}
	}

	log.Printf(
		"📧 [NOTIFIKASI EMAIL TERKIRIM] Sukses memproses pesanan #%d!\n"+
			"    ├── Penerima     : %s (User ID: %d)\n"+
			"    ├── Total Bayar  : Rp %.2f\n"+
			"    ├── Status       : PAID\n"+
			"    └── Waktu Trans. : %s",
		event.OrderID,
		recipient,
		event.UserID,
		event.TotalAmount,
		event.Timestamp.Format(time.RFC3339),
	)

	return nil
}

