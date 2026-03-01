package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	cfg := config.Load()

	assert.Equal(t, "8080", cfg.HTTPPort)
	assert.Contains(t, cfg.Postgres.DSN, "antifraude")
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.Equal(t, 10000.0, cfg.FraudRules.AmountThreshold)
	assert.Equal(t, 10, cfg.FraudRules.MaxTransactionsPerHour)
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("REDIS_DB", "2")
	t.Setenv("FRAUD_AMOUNT_THRESHOLD", "5000")
	t.Setenv("FRAUD_MAX_TX_PER_HOUR", "20")

	cfg := config.Load()

	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.Equal(t, 2, cfg.Redis.DB)
	assert.Equal(t, 5000.0, cfg.FraudRules.AmountThreshold)
	assert.Equal(t, 20, cfg.FraudRules.MaxTransactionsPerHour)
}
