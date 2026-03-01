// Package integration contains integration tests that require real infrastructure.
// Run with: go test ./tests/integration/... -tags integration
//
// Required environment variables:
//   - POSTGRES_DSN
//   - REDIS_ADDR
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/bus"
	cmdhandler "github.com/hebertzin/cqrs/internal/command/handler"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	"github.com/hebertzin/cqrs/internal/fraud/rules"
	httphandler "github.com/hebertzin/cqrs/internal/infrastructure/http/handler"
	infrahttp "github.com/hebertzin/cqrs/internal/infrastructure/http"
	"github.com/hebertzin/cqrs/internal/infrastructure/messaging/inmemory"
	pgrepository "github.com/hebertzin/cqrs/internal/infrastructure/persistence/postgres"
	redisrepository "github.com/hebertzin/cqrs/internal/infrastructure/persistence/redis"
	"github.com/hebertzin/cqrs/internal/infrastructure/projector"
	queryhandler "github.com/hebertzin/cqrs/internal/query/handler"
)

func skipIfNoInfra(t *testing.T) {
	t.Helper()
	if os.Getenv("POSTGRES_DSN") == "" {
		t.Skip("skipping integration test: POSTGRES_DSN not set")
	}
}

func buildTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	pgDSN := os.Getenv("POSTGRES_DSN")
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	pgPool, err := pgxpool.New(context.Background(), pgDSN)
	require.NoError(t, err)
	t.Cleanup(pgPool.Close)

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	t.Cleanup(func() { redisClient.Close() }) //nolint:errcheck

	log := zap.NewNop()
	eventBus := inmemory.NewEventBus(log)

	txWriteRepo := pgrepository.NewTransactionRepository(pgPool)
	accountWriteRepo := pgrepository.NewAccountRepository(pgPool)
	txReadRepo := redisrepository.NewTransactionRepository(redisClient)
	accountReadRepo := redisrepository.NewAccountRepository(redisClient)

	txProjector := projector.NewTransactionProjector(txReadRepo, accountReadRepo, log)
	txProjector.Register(eventBus)

	accountProjector := projector.NewAccountProjector(accountReadRepo, log)
	accountProjector.Register(eventBus)

	fraudEngine := rules.NewEngine(
		rules.NewAmountRule(10000),
		rules.NewLocationRule(nil),
		rules.NewVelocityRule(10, txWriteRepo),
	)

	commandBus := bus.NewCommandBus()
	commandBus.Register(cmdmodel.AnalyzeTransactionCommand, cmdhandler.NewAnalyzeTransactionHandler(
		txWriteRepo, accountWriteRepo, eventBus, fraudEngine, log,
	))
	commandBus.Register(cmdmodel.BlockAccountCommand, cmdhandler.NewBlockAccountHandler(accountWriteRepo, eventBus, log))
	commandBus.Register(cmdmodel.FlagTransactionCommand, cmdhandler.NewFlagTransactionHandler(txWriteRepo, eventBus, log))

	queryBus := bus.NewQueryBus()
	queryBus.Register(queryhandler.GetTransactionRiskKey, queryhandler.NewGetTransactionRiskHandler(txReadRepo))
	queryBus.Register(queryhandler.GetAccountStatusKey, queryhandler.NewGetAccountStatusHandler(accountReadRepo))
	queryBus.Register(queryhandler.GetFraudAlertsKey, queryhandler.NewGetFraudAlertsHandler(txReadRepo))

	router := infrahttp.NewRouter(
		httphandler.NewTransactionHandler(commandBus, queryBus, log),
		httphandler.NewAccountHandler(commandBus, queryBus, log),
		httphandler.NewHealthHandler(),
		log,
	)

	return httptest.NewServer(router)
}

func TestIntegration_AnalyzeTransaction_LowRisk(t *testing.T) {
	skipIfNoInfra(t)
	srv := buildTestServer(t)
	defer srv.Close()

	body, _ := json.Marshal(map[string]interface{}{
		"account_id":  uuid.New().String(),
		"amount":      500.00,
		"currency":    "BRL",
		"merchant_id": "merchant-safe",
		"location":    "BR",
	})

	resp, err := http.Post(srv.URL+"/api/v1/transactions", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(t, "approved", result["status"])
	assert.Equal(t, "low", result["risk_level"])
}

func TestIntegration_AnalyzeTransaction_HighRisk(t *testing.T) {
	skipIfNoInfra(t)
	srv := buildTestServer(t)
	defer srv.Close()

	body, _ := json.Marshal(map[string]interface{}{
		"account_id":  uuid.New().String(),
		"amount":      99999.99,
		"currency":    "USD",
		"merchant_id": "merchant-suspicious",
		"location":    "XX",
	})

	resp, err := http.Post(srv.URL+"/api/v1/transactions", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(t, "declined", result["status"])
	assert.Equal(t, "high", result["risk_level"])
}

func TestIntegration_HealthCheck(t *testing.T) {
	skipIfNoInfra(t)
	srv := buildTestServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
