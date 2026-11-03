package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-order-system/pkg/auth"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := auth.NewJWTManager("test-jwt-secret-key-12345", 1*time.Hour)

	validToken, err := jwtManager.Generate(123, "buyer@example.com")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectUserID   int64
	}{
		{
			name:           "No Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Authorization header format",
			authHeader:     "Basic 123456",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid JWT token",
			authHeader:     "Bearer invalid.token.payload",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid Bearer token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectUserID:   123,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(Auth(jwtManager))
			r.GET("/protected", func(c *gin.Context) {
				uid, ok := UserID(c)
				if !ok || uid != tc.expectUserID {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "user id mismatch"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d, body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

