package builder

import (
	"database/sql"
	"log"

	"go-order-system/config"
	"go-order-system/internal/http/handler"
	"go-order-system/internal/http/middleware"
	"go-order-system/internal/http/router"
	"go-order-system/internal/repository"
	"go-order-system/internal/service"
	"go-order-system/pkg/auth"
	"go-order-system/pkg/database"
	"go-order-system/pkg/route"

	"go.temporal.io/sdk/client"
)

func BuildPublicRoutes(cfg *config.Config, db *sql.DB) []route.Route {
	// 1. Inisialisasi External Clients
	rdb, err := database.NewRedisClient(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Fatalf("Gagal terhubung ke Redis: %v", err)
	}

	temporalCl, err := client.Dial(client.Options{
		HostPort: cfg.TemporalHost,
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke Temporal: %v", err)
	}

	// 2. Auth & Middleware
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)
	authMiddleware := middleware.Auth(jwtManager)

	// 3. Repositories
	userRepo := repository.NewUserRepository(db)
	rawProductRepo := repository.NewProductRepository(db)
	productRepo := repository.NewProductCacheRepository(rawProductRepo, rdb)
	orderRepo := repository.NewOrderRepository(db)

	// 4. Services
	authService := service.NewAuthService(userRepo, jwtManager)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, temporalCl)

	// 5. Handlers
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	// 6. Kembalikan kumpulan Route
	return router.PublicRoutes(authHandler, productHandler, orderHandler, authMiddleware)
}
