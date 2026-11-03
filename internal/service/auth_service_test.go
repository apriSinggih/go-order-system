package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-order-system/internal/entity"
	"go-order-system/internal/http/dto"
	"go-order-system/pkg/auth"
)

type mockUserRepo struct {
	usersByID    map[int64]*entity.User
	usersByEmail map[string]*entity.User
	nextID       int64
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		usersByID:    make(map[int64]*entity.User),
		usersByEmail: make(map[string]*entity.User),
		nextID:       1,
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error {
	user.ID = m.nextID
	m.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.usersByID[user.ID] = user
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	u, ok := m.usersByEmail[email]
	if !ok {
		return nil, fmt.Errorf("pengguna tidak ditemukan")
	}
	return u, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, fmt.Errorf("pengguna tidak ditemukan")
	}
	return u, nil
}

func TestAuthService_RegisterAndLogin(t *testing.T) {
	repo := newMockUserRepo()
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour)
	authSvc := NewAuthService(repo, jwtManager)
	ctx := context.Background()

	// 1. Register Success
	regReq := dto.RegisterRequest{
		Email:    "buyer@example.com",
		Password: "password123",
	}
	user, token, err := authSvc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}
	if user.Email != "buyer@example.com" {
		t.Errorf("expected email buyer@example.com, got %s", user.Email)
	}
	if token == "" {
		t.Errorf("expected non-empty token")
	}

	// 2. Register Duplicate Email Fails
	_, _, err = authSvc.Register(ctx, regReq)
	if err == nil {
		t.Fatalf("expected error for duplicate email registration, got nil")
	}

	// 3. Login Success
	loginReq := dto.LoginRequest{
		Email:    "buyer@example.com",
		Password: "password123",
	}
	loginToken, err := authSvc.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}
	if loginToken == "" {
		t.Errorf("expected non-empty login token")
	}

	// 4. Login Wrong Password
	_, err = authSvc.Login(ctx, dto.LoginRequest{
		Email:    "buyer@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatalf("expected error for wrong password, got nil")
	}

	// 5. Login Unknown Email
	_, err = authSvc.Login(ctx, dto.LoginRequest{
		Email:    "unknown@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatalf("expected error for unknown email, got nil")
	}

	// 6. GetProfile Success
	profile, err := authSvc.GetProfile(ctx, user.ID)
	if err != nil {
		t.Fatalf("unexpected error getting profile: %v", err)
	}
	if profile.ID != user.ID || profile.Email != user.Email {
		t.Errorf("profile mismatch: got %+v, want %+v", profile, user)
	}

	// 7. GetProfile Unknown User
	_, err = authSvc.GetProfile(ctx, 9999)
	if err == nil {
		t.Fatalf("expected error for unknown user profile, got nil")
	}
}

