package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// CORSMiddleware mengizinkan request dari frontend origin (misal localhost:5173)
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func RateLimiter(rdb *redis.Client, limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		ctx := c.Request.Context()
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			c.Abort()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		if count > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak permintaan. Silakan coba lagi nanti.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}