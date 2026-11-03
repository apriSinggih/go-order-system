package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-order-system/internal/entity"
	"go-order-system/internal/http/dto"

	"github.com/gin-gonic/gin"
)

type mockOrderService struct {
	createOrderFn func(ctx context.Context, userID int64, req dto.CreateOrderRequest) (*entity.Order, error)
	getOrderFn    func(ctx context.Context, userID, id int64) (*entity.Order, error)
}

func (m *mockOrderService) CreateOrder(ctx context.Context, userID int64, req dto.CreateOrderRequest) (*entity.Order, error) {
	if m.createOrderFn != nil {
		return m.createOrderFn(ctx, userID, req)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockOrderService) GetOrder(ctx context.Context, userID, id int64) (*entity.Order, error) {
	if m.getOrderFn != nil {
		return m.getOrderFn(ctx, userID, id)
	}
	return nil, fmt.Errorf("not implemented")
}

func TestOrderHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockOrderService{
		createOrderFn: func(ctx context.Context, userID int64, req dto.CreateOrderRequest) (*entity.Order, error) {
			if len(req.Items) == 0 {
				return nil, fmt.Errorf("item tidak boleh kosong")
			}
			return &entity.Order{
				ID:          1,
				UserID:      userID,
				TotalAmount: 500000,
				Status:      "PENDING",
			}, nil
		},
	}

	h := NewOrderHandler(mockSvc)

	// 1. Unauthenticated (no user_id in context)
	rUnauth := gin.New()
	rUnauth.POST("/orders", h.Create)

	body, _ := json.Marshal(dto.CreateOrderRequest{
		Items: []dto.CreateOrderItemRequest{{ProductID: 1, Quantity: 2}},
	})
	req1, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	rUnauth.ServeHTTP(w1, req1)

	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", w1.Code)
	}

	// 2. Authenticated + Valid Request
	rAuth := gin.New()
	rAuth.Use(func(c *gin.Context) {
		c.Set("auth_user_id", int64(42))
		c.Next()
	})
	rAuth.POST("/orders", h.Create)

	req2, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	rAuth.ServeHTTP(w2, req2)

	if w2.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w2.Code)
	}
}

func TestOrderHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockOrderService{
		getOrderFn: func(ctx context.Context, userID, id int64) (*entity.Order, error) {
			if userID == 42 && id == 10 {
				return &entity.Order{ID: 10, UserID: 42, Status: "PAID"}, nil
			}
			return nil, fmt.Errorf("pesanan tidak ditemukan")
		},
	}

	h := NewOrderHandler(mockSvc)

	// 1. Unauthenticated
	rUnauth := gin.New()
	rUnauth.GET("/orders/:id", h.GetByID)

	req1, _ := http.NewRequest(http.MethodGet, "/orders/10", nil)
	w1 := httptest.NewRecorder()
	rUnauth.ServeHTTP(w1, req1)

	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", w1.Code)
	}

	// 2. Authenticated + Found
	rAuth := gin.New()
	rAuth.Use(func(c *gin.Context) {
		c.Set("auth_user_id", int64(42))
		c.Next()
	})
	rAuth.GET("/orders/:id", h.GetByID)

	req2, _ := http.NewRequest(http.MethodGet, "/orders/10", nil)
	w2 := httptest.NewRecorder()
	rAuth.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w2.Code)
	}

	// 3. Authenticated + Not Found (other user's order)
	req3, _ := http.NewRequest(http.MethodGet, "/orders/99", nil)
	w3 := httptest.NewRecorder()
	rAuth.ServeHTTP(w3, req3)

	if w3.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", w3.Code)
	}

	// 4. Invalid ID format
	req4, _ := http.NewRequest(http.MethodGet, "/orders/abc", nil)
	w4 := httptest.NewRecorder()
	rAuth.ServeHTTP(w4, req4)

	if w4.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w4.Code)
	}
}

