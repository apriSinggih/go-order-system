package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"go-order-system/internal/repository"
	httpHandler "go-order-system/internal/transport/http"
	"go-order-system/internal/usecase"
	"go-order-system/pkg/database"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "secretpassword")
	dbName := getEnv("DB_NAME", "order_db")

	redisHost := getEnv("REDIS_HOST", "redis")
	redisPort := getEnv("REDIS_PORT", "6379")

	temporalHost := getEnv("TEMPORAL_HOST", "temporal:7233")

	// 1. Inisialisasi Database PostgreSQL & Redis
	dbConfig := database.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPass,
		DBName:   dbName,
		SSLMode:  "disable",
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	rdb, err := database.NewRedisClient(redisHost, redisPort)
	if err != nil {
		log.Fatalf("Gagal terhubung ke Redis: %v", err)
	}

	// 2. Inisialisasi Temporal Client
	temporalCl, err := client.Dial(client.Options{
		HostPort: temporalHost,
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke Temporal Server: %v", err)
	}
	defer temporalCl.Close()

	// 3. Auto Migration
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	migrationsPath := getEnv("MIGRATIONS_PATH", "deployments/migrations")
	if err := database.RunMigrations(migrationsPath, dbURL); err != nil {
		log.Fatalf("Gagal menjalankan migrasi database: %v", err)
	}

	// 4. Repositories
	rawProductRepo := repository.NewProductRepository(db)
	productRepo := repository.NewProductCacheRepository(rawProductRepo, rdb)
	orderRepo := repository.NewOrderRepository(db)

	// 5. Usecases
	productUsecase := usecase.NewProductUsecase(productRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, productRepo, temporalCl)

	// 6. Router & Handlers
	r := gin.Default()
	httpHandler.InitMetrics(r)
	r.Use(httpHandler.CORSMiddleware())
	r.Use(httpHandler.RateLimiter(rdb, 60, time.Minute))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "Go API connected to Postgres, Redis, & Temporal!"})
	})

	httpHandler.NewProductHandler(r, productUsecase)
	httpHandler.NewOrderHandler(r, orderUsecase)

	port := getEnv("PORT", "8080")
	fmt.Printf("Server HTTP berjalan di port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}