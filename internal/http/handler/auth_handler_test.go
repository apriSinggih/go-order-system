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

type mockAuthService struct {
	registerFn   func(ctx context.Context, req dto.RegisterRequest) (*entity.User, string, error)
	loginFn      func(ctx context.Context, req dto.LoginRequest) (string, error)
	getProfileFn func(ctx context.Context, userID int64) (*entity.User, error)
}

func (m *mockAuthService) Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, string, error) {
	if m.registerFn != nil {
		return m.registerFn(ctx, req)
	}
	return nil, "", fmt.Errorf("not implemented")
}

func (m *mockAuthService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, req)
	}
	return "", fmt.Errorf("not implemented")
}

func (m *mockAuthService) GetProfile(ctx context.Context, userID int64) (*entity.User, error) {
	if m.getProfileFn != nil {
		return m.getProfileFn(ctx, userID)
	}
	return nil, fmt.Errorf("not implemented")
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockAuthService{
		registerFn: func(ctx context.Context, req dto.RegisterRequest) (*entity.User, string, error) {
			if req.Email == "duplicate@example.com" {
				return nil, "", fmt.Errorf("email sudah terdaftar")
			}
			return &entity.User{ID: 1, Email: req.Email}, "jwt-token-123", nil
		},
	}

	h := NewAuthHandler(mockSvc)
	r := gin.New()
	r.POST("/register", h.Register)

	// Valid registration
	body, _ := json.Marshal(dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
	}

	// Duplicate registration
	bodyDup, _ := json.Marshal(dto.RegisterRequest{
		Email:    "duplicate@example.com",
		Password: "password123",
	})
	reqDup, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyDup))
	reqDup.Header.Set("Content-Type", "application/json")
	wDup := httptest.NewRecorder()
	r.ServeHTTP(wDup, reqDup)

	if wDup.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, wDup.Code)
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockAuthService{
		loginFn: func(ctx context.Context, req dto.LoginRequest) (string, error) {
			if req.Password != "correctpass" {
				return "", fmt.Errorf("email atau password salah")
			}
			return "valid-token", nil
		},
	}

	h := NewAuthHandler(mockSvc)
	r := gin.New()
	r.POST("/login", h.Login)

	// Success
	body, _ := json.Marshal(dto.LoginRequest{
		Email:    "test@example.com",
		Password: "correctpass",
	})
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}

	// Unauthorized
	bodyFail, _ := json.Marshal(dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpass",
	})
	reqFail, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyFail))
	reqFail.Header.Set("Content-Type", "application/json")
	wFail := httptest.NewRecorder()
	r.ServeHTTP(wFail, reqFail)

	if wFail.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, wFail.Code)
	}
}

func TestAuthHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockAuthService{
		getProfileFn: func(ctx context.Context, userID int64) (*entity.User, error) {
			if userID == 1 {
				return &entity.User{ID: 1, Email: "user@example.com"}, nil
			}
			return nil, fmt.Errorf("pengguna tidak ditemukan")
		},
	}

	h := NewAuthHandler(mockSvc)

	// Case 1: Unauthenticated (no user_id in context)
	rUnauth := gin.New()
	rUnauth.GET("/me", h.Me)

	req1, _ := http.NewRequest(http.MethodGet, "/me", nil)
	w1 := httptest.NewRecorder()
	rUnauth.ServeHTTP(w1, req1)

	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", w1.Code)
	}

	// Case 2: Authenticated and found
	rAuth := gin.New()
	rAuth.Use(func(c *gin.Context) {
		c.Set("auth_user_id", int64(1))
		c.Next()
	})
	rAuth.GET("/me", h.Me)

	req2, _ := http.NewRequest(http.MethodGet, "/me", nil)
	w2 := httptest.NewRecorder()
	rAuth.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w2.Code)
	}
}


