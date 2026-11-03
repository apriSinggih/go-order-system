
# Go Order Processing & Distributed Transaction System

Sistem backend *production-grade* berbasis Microservices & Clean Architecture menggunakan **Go**, **PostgreSQL**, **Redis**, **Temporal Workflow (Saga Pattern)**, dan **Apache Kafka** untuk memproses transaksi checkout pesanan secara terdistribusi.

---

## 🚀 Fitur Utama

- **Clean Architecture & SOLID Principles**: Pemisahan layer yang jelas (`domain`, `usecase`, `repository`, `transport`).
- **Fully Containerized Environment**: Seluruh *stack* infrastruktur dan aplikasi berjalan terisolasi menggunakan **Docker Desktop & Docker Compose**.
- **Automated Database Migration**: Pengelolaan skema PostgreSQL terprogram menggunakan `golang-migrate` secara otomatis saat aplikasi dinyalakan.
- **Caching & Performance**: Menggunakan **Redis** sebagai *caching layer* (Decorator Pattern) dan *HTTP Rate Limiter Middleware*.
- **Durable Orchestration (Saga Pattern)**: Menggunakan **Temporal.io** untuk menangani transaksi terdistribusi (*Deduct Stock*, *Charge Payment*, *Auto-Rollback / Release Stock* jika pembayaran gagal).
- **Event-Driven Architecture (EDA)**: *Publishing* event `order.paid` secara *asynchronous* ke **Apache Kafka** setelah transaksi berhasil diproses.

---

## 🛠️ Tech Stack & Infrastructure

- **Language**: Go (Golang)
- **Web Framework**: Gin Gonic
- **Database**: PostgreSQL (via SQLX)
- **Cache & Rate Limiter**: Redis
- **Workflow Engine**: Temporal.io
- **Event Streaming / Message Broker**: Apache Kafka & Zookeeper
- **Migration Tool**: `golang-migrate`
- **Containerization**: Docker & Docker Compose

---

## 📁 Struktur Direktori

```text
.
├── cmd
│   └── app
│       └── main.go
├── config
│   └── config.go
├── db
│   └── migrations
├── go.mod
├── go.sum
├── internal
│   ├── builder
│   │   └── builder.go
│   ├── entity
│   │   ├── jadwal.go
│   │   ├── kursi.go
│   │   ├── pemesanan.go
│   │   ├── pengguna.go
│   │   └── rute.go
│   ├── http
│   │   ├── dto
│   │   │   └── pemesanan_dto.go
│   │   ├── handler
│   │   │   ├── pemesanan_handler.go
│   │   │   └── seed_handler.go
│   │   └── router
│   │       └── router.go
│   ├── repository
│   │   ├── jadwal_repository.go
│   │   ├── kursi_repository.go
│   │   └── pemesanan_repository.go
│   └── service
│       └── pemesanan_service.go
├── pkg
│   ├── database
│   │   └── postgres.go
│   ├── response
│   │   └── response.go
│   ├── route
│   │   └── route.go
│   └── server
│       └── echo.go