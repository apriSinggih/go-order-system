package main

import (
	"log"

	"go-order-system/config"
	"go-order-system/internal/repository"
	orderWorkflow "go-order-system/internal/workflow"
	"go-order-system/pkg/database"
	"go-order-system/pkg/messaging"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	cfg, err := config.NewConfig(".env")
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// 1. Inisialisasi Database PostgreSQL
	db, err := database.InitDatabase(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Worker gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	// 2. Inisialisasi Kafka Producer
	kafkaProducer, err := messaging.NewKafkaProducer(cfg.KafkaBroker)
	if err != nil {
		log.Fatalf("Worker gagal terhubung ke Kafka: %v", err)
	}
	defer kafkaProducer.Close()

	// 3. Inisialisasi Temporal Client
	temporalClient, err := client.Dial(client.Options{
		HostPort: cfg.TemporalHost,
	})
	if err != nil {
		log.Fatalf("Worker gagal terhubung ke Temporal Server: %v", err)
	}
	defer temporalClient.Close()

	// 4. Inisialisasi Repositories & Activities
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	activities := orderWorkflow.NewOrderActivities(productRepo, orderRepo, kafkaProducer)

	// 5. Mendaftarkan & Menjalankan Worker
	w := worker.New(temporalClient, orderWorkflow.OrderTaskQueue, worker.Options{})
	w.RegisterWorkflow(orderWorkflow.ProcessOrderWorkflow)
	w.RegisterActivity(activities)

	log.Println("⚡ Temporal Worker siap mendengarkan tugas pada queue:", orderWorkflow.OrderTaskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("Gagal menjalankan Temporal worker: %v", err)
	}
}