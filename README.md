
# Go Order Processing & Distributed Transaction System

Sistem backend *production-grade* berbasis Microservices & Clean Architecture menggunakan **Go**, **PostgreSQL**, **Redis**, **Temporal Workflow (Saga Pattern)**, **Apache Kafka**, dan **Prometheus & Grafana** untuk memproses transaksi pesanan e-commerce secara terdistribusi.

---

## 🚀 Fitur Utama

- **Clean Architecture & SOLID Principles**: Pemisahan layer yang jelas (`entity`, `repository`, `service`, `handler`, `builder`).
- **Fully Containerized Environment**: Seluruh *stack* infrastruktur dan microservices berjalan terisolasi menggunakan **Docker & Docker Compose**.
- **Automated Database Migration**: Pengelolaan skema PostgreSQL terprogram menggunakan `golang-migrate` otomatis saat aplikasi dinyalakan.
- **Caching & High Performance**: Menggunakan **Redis** sebagai *Read-Through Cache* (Decorator Pattern) untuk produk katalog.
- **Authentication & Authorization (JWT)**: Registrasi & login pengguna berbasis Bcrypt hash dan JWT Bearer authentication token.
- **Durable Orchestration (Saga Pattern)**: Menggunakan **Temporal.io** untuk menangani transaksi terdistribusi (*Deduct Stock*, *Charge Payment*, *Auto-Rollback / Release Stock* jika pembayaran gagal, *Update Order Status*).
- **Event-Driven Architecture (EDA)**: *Publishing* event `order.paid` ke **Apache Kafka** setelah transaksi sukses.
- **Notification Consumer Microservice**: Worker mandiri yang mendengarkan event dari Kafka dan memproses notifikasi / invoice ke pembeli.
- **Observability & Monitoring**: Prometheus scraping metrik HTTP `/metrics` dan dashboard Grafana ter-provisioning otomatis.

---

## 🛠️ Tech Stack & Infrastructure

- **Language**: Go (Golang)
- **Web Framework**: Gin Gonic
- **Database**: PostgreSQL 15 (`database/sql`)
- **Cache**: Redis 7
- **Distributed Workflow**: Temporal.io
- **Message Broker**: Apache Kafka & Zookeeper
- **Migration Tool**: `golang-migrate`
- **Observability**: Prometheus & Grafana
- **Containerization**: Docker & Docker Compose

---

## 📁 Struktur Direktori

```text
.
├── cmd/
│   ├── app/          # Entrypoint REST API
│   ├── consumer/     # Entrypoint Kafka Consumer (Notification Worker)
│   └── worker/       # Entrypoint Temporal Worker & Kafka Producer
├── config/           # Konfigurasi & Environment (.env)
├── database/
│   └── migrations/   # Skema migrasi SQL (.up.sql & .down.sql)
├── deployments/
│   ├── grafana/      # Konfigurasi datasource & dashboard Grafana
│   └── prometheus/   # Konfigurasi Prometheus scraper
├── internal/
│   ├── builder/      # Dependency Injection Builder
│   ├── entity/       # Domain Structs (Product, Order, User, Events)
│   ├── http/
│   │   ├── dto/      # Request & Response DTOs
│   │   ├── handler/  # HTTP Handlers (Auth, Product, Order)
│   │   ├── middleware/ # JWT Auth & Prometheus Metrics Middleware
│   │   └── router/   # Definisi route API
│   ├── repository/   # Data Access Layer (Postgres & Redis Cache)
│   ├── service/      # Business Logic Services
│   └── workflow/     # Temporal Saga Workflow & Activities
└── pkg/
    ├── auth/         # JWT Token Manager
    ├── database/     # DB Connectors & Auto-migration
    ├── messaging/    # Kafka Producer & Consumer
    ├── response/     # Standar respon JSON API
    ├── route/        # Abstraksi Route & Middleware
    └── server/       # HTTP Server dengan Graceful Shutdown
```

---

## 📡 Daftar Endpoint API

### Autentikasi
- `POST /api/v1/auth/register` - Pendaftaran user baru
- `POST /api/v1/auth/login` - Login user & perolehan JWT Bearer Token
- `GET /api/v1/auth/me` - Mengambil profil user yang sedang login (Memerlukan `Authorization: Bearer <token>`)

### Produk (Cached via Redis)
- `GET /api/v1/products` - Daftar semua produk
- `GET /api/v1/products/:id` - Detail produk berdasarkan ID

### Pemesanan (Memerlukan `Authorization: Bearer <token>`)
- `POST /api/v1/orders` - Membuat order baru (memicu Temporal Saga Workflow)
- `GET /api/v1/orders/:id` - Mengambil detail pesanan milik pengguna

### Sistem & Monitoring
- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics scraper endpoint (port 8080)
- `http://localhost:3000` - Dashboard Grafana (admin / admin)
- `http://localhost:8233` - Temporal Web UI
