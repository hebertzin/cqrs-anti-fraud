package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/bus"
	cmdhandler "github.com/hebertzin/cqrs/internal/command/handler"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	"github.com/hebertzin/cqrs/internal/config"
	"github.com/hebertzin/cqrs/internal/fraud/rules"
	infrahttp "github.com/hebertzin/cqrs/internal/infrastructure/http"
	httphandler "github.com/hebertzin/cqrs/internal/infrastructure/http/handler"
	"github.com/hebertzin/cqrs/internal/infrastructure/messaging/inmemory"
	pgrepository "github.com/hebertzin/cqrs/internal/infrastructure/persistence/postgres"
	redisrepository "github.com/hebertzin/cqrs/internal/infrastructure/persistence/redis"
	"github.com/hebertzin/cqrs/internal/infrastructure/projector"
	queryhandler "github.com/hebertzin/cqrs/internal/query/handler"
	"github.com/hebertzin/cqrs/pkg/logger"
)

func main() {
	cfg := config.Load()

	log, err := logger.New(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer log.Sync() //nolint:errcheck

	pgPool, err := pgxpool.New(context.Background(), cfg.Postgres.DSN)
	if err != nil {
		log.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer pgPool.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close() //nolint:errcheck

	txWriteRepo := pgrepository.NewTransactionRepository(pgPool)
	accountWriteRepo := pgrepository.NewAccountRepository(pgPool)
	txReadRepo := redisrepository.NewTransactionRepository(redisClient)
	accountReadRepo := redisrepository.NewAccountRepository(redisClient)

	eventBus := inmemory.NewEventBus(log)

	txProjector := projector.NewTransactionProjector(txReadRepo, accountReadRepo, log)
	txProjector.Register(eventBus)

	accountProjector := projector.NewAccountProjector(accountReadRepo, log)
	accountProjector.Register(eventBus)

	fraudEngine := rules.NewEngine(
		rules.NewAmountRule(cfg.FraudRules.AmountThreshold),
		rules.NewVelocityRule(cfg.FraudRules.MaxTransactionsPerHour, txWriteRepo),
		rules.NewLocationRule(nil),
	)

	commandBus := bus.NewCommandBus()
	commandBus.Register(cmdmodel.AnalyzeTransactionCommand, cmdhandler.NewAnalyzeTransactionHandler(
		txWriteRepo, accountWriteRepo, eventBus, fraudEngine, log,
	))
	commandBus.Register(cmdmodel.BlockAccountCommand, cmdhandler.NewBlockAccountHandler(
		accountWriteRepo, eventBus, log,
	))
	commandBus.Register(cmdmodel.FlagTransactionCommand, cmdhandler.NewFlagTransactionHandler(
		txWriteRepo, eventBus, log,
	))

	queryBus := bus.NewQueryBus()
	queryBus.Register(queryhandler.GetTransactionRiskKey, queryhandler.NewGetTransactionRiskHandler(txReadRepo))
	queryBus.Register(queryhandler.GetAccountStatusKey, queryhandler.NewGetAccountStatusHandler(accountReadRepo))
	queryBus.Register(queryhandler.GetFraudAlertsKey, queryhandler.NewGetFraudAlertsHandler(txReadRepo))

	txHTTPHandler := httphandler.NewTransactionHandler(commandBus, queryBus, log)
	accountHTTPHandler := httphandler.NewAccountHandler(commandBus, queryBus, log)
	healthHTTPHandler := httphandler.NewHealthHandler()

	router := infrahttp.NewRouter(txHTTPHandler, accountHTTPHandler, healthHTTPHandler, log)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced shutdown", zap.Error(err))
	}
	log.Info("server stopped")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
