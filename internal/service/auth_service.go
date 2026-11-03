package service

import (
	"context"
	"fmt"
	"strings"

	"go-order-system/internal/entity"
	"go-order-system/internal/http/dto"
	"go-order-system/internal/repository"
	"go-order-system/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, string, error)
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userID int64) (*entity.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwt      *auth.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwt *auth.JWTManager) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, "", fmt.Errorf("email sudah terdaftar")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	user := &entity.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := s.jwt.Generate(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("email atau password salah")
	}

	return s.jwt.Generate(user.ID, user.Email)
}

func (s *authService) GetProfile(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("pengguna tidak ditemukan: %w", err)
	}
	return user, nil
}

