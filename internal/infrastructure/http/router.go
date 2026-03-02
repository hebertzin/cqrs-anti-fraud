package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/hebertzin/cqrs/internal/infrastructure/http/handler"
	"github.com/hebertzin/cqrs/internal/infrastructure/http/middleware"
	"go.uber.org/zap"
)

func NewRouter(
	transactionHandler *handler.TransactionHandler,
	accountHandler *handler.AccountHandler,
	healthHandler *handler.HealthHandler,
	logger *zap.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Recover(logger))
	r.Use(middleware.Logger(logger))

	r.Get("/health", healthHandler.Live)
	r.Get("/ready", healthHandler.Ready)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/transactions", func(r chi.Router) {
			r.Post("/", transactionHandler.AnalyzeTransaction)
			r.Get("/{id}/risk", transactionHandler.GetTransactionRisk)
			r.Post("/{id}/flag", transactionHandler.FlagTransaction)
		})

		r.Route("/accounts", func(r chi.Router) {
			r.Get("/{id}/status", accountHandler.GetAccountStatus)
			r.Post("/{id}/block", accountHandler.BlockAccount)
		})

		r.Route("/fraud", func(r chi.Router) {
			r.Get("/alerts", accountHandler.GetFraudAlerts)
		})
	})

	return r
}
