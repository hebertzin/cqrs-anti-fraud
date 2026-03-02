package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort   string
	Postgres   PostgresConfig
	Redis      RedisConfig
	FraudRules FraudRulesConfig
}

type PostgresConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type FraudRulesConfig struct {
	AmountThreshold        float64
	MaxTransactionsPerHour int
}

func Load() Config {
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	amountThreshold, _ := strconv.ParseFloat(getEnv("FRAUD_AMOUNT_THRESHOLD", "10000"), 64)
	maxTxPerHour, _ := strconv.Atoi(getEnv("FRAUD_MAX_TX_PER_HOUR", "10"))

	return Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		Postgres: PostgresConfig{
			DSN: getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/antifraude?sslmode=disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		FraudRules: FraudRulesConfig{
			AmountThreshold:        amountThreshold,
			MaxTransactionsPerHour: maxTxPerHour,
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
