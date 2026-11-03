package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	DBUrl        string
	RedisHost    string
	RedisPort    string
	TemporalHost string
	KafkaBroker  string
	JWTSecret    string
	JWTExpiry    time.Duration
}

func NewConfig(envFile string) (*Config, error) {
	_ = godotenv.Load(envFile)

	expiryStr := getEnv("JWT_EXPIRY", "24h")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		expiry = 24 * time.Hour
	}

	return &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		DBUrl:        getEnv("DB_URL", "postgres://postgres:secretpassword@postgres:5432/order_db?sslmode=disable"),
		RedisHost:    getEnv("REDIS_HOST", "redis"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		TemporalHost: getEnv("TEMPORAL_HOST", "temporal:7233"),
		KafkaBroker:  getEnv("KAFKA_BROKER", "kafka:9092"),
		JWTSecret:    getEnv("JWT_SECRET", "super-secret-jwt-key-for-go-order-system"),
		JWTExpiry:    expiry,
	}, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
