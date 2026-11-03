package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-order-system/config"
	"go-order-system/internal/entity"
	"go-order-system/internal/repository"
	"go-order-system/internal/service"
	"go-order-system/pkg/database"
	"go-order-system/pkg/messaging"
)

func main() {
	cfg, err := config.NewConfig(".env")
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// 1. Inisialisasi Database PostgreSQL (untuk mengambil data profil pengguna)
	db, err := database.InitDatabase(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Consumer gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	notificationService := service.NewNotificationService(userRepo)

	// 2. Inisialisasi Kafka Consumer
	const groupID = "notification-consumer-group"
	consumer := messaging.NewKafkaConsumer(cfg.KafkaBroker, groupID, messaging.TopicOrderPaid)
	defer consumer.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("📬 Notification Consumer Service berjalan, mendengarkan topic '%s'...", messaging.TopicOrderPaid)

	for {
		msg, err := consumer.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Consumer menerima sinyal penghentian, mematikan consumer...")
				break
			}
			log.Printf("Gagal membaca pesan dari Kafka: %v", err)
			continue
		}

		var event entity.OrderPaidEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Gagal mengurai event OrderPaidEvent: %v", err)
			continue
		}

		if err := notificationService.SendOrderPaidNotification(ctx, event); err != nil {
			log.Printf("Gagal mengirim notifikasi pesanan #%d: %v", event.OrderID, err)
		}
	}

	log.Println("Notification Consumer Service berhasil dimatikan dengan aman.")
}

