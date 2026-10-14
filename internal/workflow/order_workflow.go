package workflow

import (
	"context"
	"fmt"
	"time"

	"go-order-system/internal/entity"
	"go-order-system/internal/repository"
	"go-order-system/pkg/messaging"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const OrderTaskQueue = "ORDER_TASK_QUEUE"

type OrderActivities struct {
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
	kafka       *messaging.KafkaProducer
}

func NewOrderActivities(
	productRepo repository.ProductRepository,
	orderRepo repository.OrderRepository,
	kafka *messaging.KafkaProducer,
) *OrderActivities {
	return &OrderActivities{
		productRepo: productRepo,
		orderRepo:   orderRepo,
		kafka:       kafka,
	}
}

// ----------------- ACTIVITIES -----------------

func (a *OrderActivities) DeductStockActivity(ctx context.Context, order *entity.Order) error {
	for _, item := range order.Items {
		err := a.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return fmt.Errorf("gagal memotong stok produk %d: %w", item.ProductID, err)
		}
	}
	return nil
}

func (a *OrderActivities) ReleaseStockActivity(ctx context.Context, order *entity.Order) error {
	for _, item := range order.Items {
		_ = a.productRepo.UpdateStock(ctx, item.ProductID, -item.Quantity)
	}
	return nil
}

func (a *OrderActivities) ChargePaymentActivity(ctx context.Context, order *entity.Order) error {
	// Simulasi limit payment 50 juta
	if order.TotalAmount > 50000000 {
		return fmt.Errorf("pembayaran ditolak: total transaksi melebihi limit Rp50.000.000")
	}
	return nil
}

func (a *OrderActivities) UpdateOrderStatusActivity(ctx context.Context, orderID int64, status string) error {
	return a.orderRepo.UpdateStatus(ctx, orderID, status)
}

func (a *OrderActivities) PublishOrderPaidEventActivity(ctx context.Context, order *entity.Order) error {
	event := entity.OrderPaidEvent{
		OrderID:     order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Timestamp:   time.Now(),
	}

	key := fmt.Sprintf("order-%d", order.ID)
	return a.kafka.PublishEvent(ctx, key, event)
}

// ----------------- WORKFLOW (SAGA PATTERN) -----------------

func ProcessOrderWorkflow(ctx workflow.Context, order *entity.Order) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var activities *OrderActivities

	// 1. Deduct Stock
	err := workflow.ExecuteActivity(ctx, activities.DeductStockActivity, order).Get(ctx, nil)
	if err != nil {
		_ = workflow.ExecuteActivity(ctx, activities.UpdateOrderStatusActivity, order.ID, "FAILED").Get(ctx, nil)
		return err
	}

	// 2. Charge Payment
	err = workflow.ExecuteActivity(ctx, activities.ChargePaymentActivity, order).Get(ctx, nil)
	if err != nil {
		// SAGA COMPENSATION: Kembalikan Stok jika Payment Gagal
		_ = workflow.ExecuteActivity(ctx, activities.ReleaseStockActivity, order).Get(ctx, nil)
		_ = workflow.ExecuteActivity(ctx, activities.UpdateOrderStatusActivity, order.ID, "FAILED").Get(ctx, nil)
		return fmt.Errorf("pembayaran gagal, stok dikembalikan: %w", err)
	}

	// 3. Update Status ke PAID
	err = workflow.ExecuteActivity(ctx, activities.UpdateOrderStatusActivity, order.ID, "PAID").Get(ctx, nil)
	if err != nil {
		return err
	}

	// 4. Publish Event ke Kafka
	_ = workflow.ExecuteActivity(ctx, activities.PublishOrderPaidEventActivity, order).Get(ctx, nil)

	return nil
}