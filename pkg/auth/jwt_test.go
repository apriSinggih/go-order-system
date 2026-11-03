package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTManager_GenerateAndParse(t *testing.T) {
	manager := NewJWTManager("test-secret-key-123", 1*time.Hour)

	userID := int64(42)
	email := "test@example.com"

	token, err := manager.Generate(userID, email)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("unexpected error parsing token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-123", 1*time.Hour)

	now := time.Now().Add(-2 * time.Hour)
	claims := Claims{
		UserID: 1,
		Email:  "expired@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("test-secret-key-123"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = manager.Parse(tokenStr)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}
}

func TestJWTManager_InvalidSecret(t *testing.T) {
	manager1 := NewJWTManager("secret-1", 1*time.Hour)
	manager2 := NewJWTManager("secret-2", 1*time.Hour)

	token, err := manager1.Generate(1, "user@example.com")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = manager2.Parse(token)
	if err == nil {
		t.Fatalf("expected error when validating with wrong secret, got nil")
	}
}
