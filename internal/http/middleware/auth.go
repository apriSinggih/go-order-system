package middleware

import (
	"net/http"
	"strings"

	"go-order-system/pkg/auth"
	"go-order-system/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	userIDKey = "auth_user_id"
	emailKey  = "auth_email"
)

// Auth memvalidasi Bearer token JWT dan menyimpan identitas pengguna di context.
func Auth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "Token otentikasi diperlukan")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims, err := jwtManager.Parse(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set(userIDKey, claims.UserID)
		c.Set(emailKey, claims.Email)
		c.Next()
	}
}

// UserID mengambil ID pengguna yang sudah divalidasi dari context.
func UserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(userIDKey)
	if !exists {
		return 0, false
	}
	userID, ok := value.(int64)
	return userID, ok
}

// Email mengambil email pengguna dari context.
func Email(c *gin.Context) (string, bool) {
	value, exists := c.Get(emailKey)
	if !exists {
		return "", false
	}
	email, ok := value.(string)
	return email, ok
}
