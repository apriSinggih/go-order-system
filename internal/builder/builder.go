package builder

import (
	"database/sql"
	"log"

	"go-order-system/config"
	"go-order-system/internal/http/handler"
	"go-order-system/internal/http/router"
	"go-order-system/internal/repository"
	"go-order-system/internal/service"
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

	// 2. Repositories
	rawProductRepo := repository.NewProductRepository(db)
	productRepo := repository.NewProductCacheRepository(rawProductRepo, rdb)
	orderRepo := repository.NewOrderRepository(db)

	// 3. Services
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, temporalCl)

	// 4. Handlers
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	// 5. Kembalikan kumpulan Route
	return router.PublicRoutes(productHandler, orderHandler)
}