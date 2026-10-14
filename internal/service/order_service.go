package service

import (
	"context"
	"fmt"

	"go-order-system/internal/entity"
	"go-order-system/internal/http/dto"
	"go-order-system/internal/repository"

	"go.temporal.io/sdk/client"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*entity.Order, error)
	GetOrder(ctx context.Context, id int64) (*entity.Order, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	temporalCl  client.Client
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	temporalCl client.Client,
) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		temporalCl:  temporalCl,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*entity.Order, error) {
	var totalAmount float64
	var orderItems []entity.OrderItem

	for _, item := range req.Items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("produk ID %d tidak ditemukan: %w", item.ProductID, err)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("stok produk '%s' tidak mencukupi (sisa: %d)", product.Name, product.Stock)
		}

		totalAmount += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	order := &entity.Order{
		UserID:      req.UserID,
		TotalAmount: totalAmount,
		Status:      "PENDING",
	}

	if err := s.orderRepo.Create(ctx, order, orderItems); err != nil {
		return nil, fmt.Errorf("gagal membuat pesanan: %w", err)
	}
	order.Items = orderItems

	// Trigger Temporal Workflow (Proses Async Saga)
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("order-workflow-%d", order.ID),
		TaskQueue: "ORDER_TASK_QUEUE",
	}

	_, err := s.temporalCl.ExecuteWorkflow(ctx, workflowOptions, "ProcessOrderWorkflow", order)
	if err != nil {
		return nil, fmt.Errorf("gagal memulai workflow transaksi: %w", err)
	}

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, id int64) (*entity.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}