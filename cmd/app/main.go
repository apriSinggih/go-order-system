package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"go-order-system/config"
	"go-order-system/internal/builder"
	"go-order-system/pkg/database"
	"go-order-system/pkg/server"
)

func main() {
	cfg, err := config.NewConfig(".env")
	checkError(err)

	db, err := database.InitDatabase(cfg.DBUrl)
	checkError(err)
	defer db.Close()

	publicRoutes := builder.BuildPublicRoutes(cfg, db)

	srv := server.NewServer(cfg, publicRoutes)

	fmt.Printf("🚀 Server is running on port %s\n", cfg.AppPort)
	runServer(srv, cfg.AppPort)
	waitForShutdown(srv)
}

func waitForShutdown(srv *server.Server) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	<-quit
	fmt.Println("\nMematikan server secara perlahan (Graceful Shutdown)...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			srv.Logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Server berhasil dimatikan")
}

func runServer(srv *server.Server, port string) {
	go func() {
		err := srv.Start(fmt.Sprintf(":%s", port))
		if err != nil && err.Error() != "http: Server closed" {
			log.Fatal(err)
		}
	}()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}