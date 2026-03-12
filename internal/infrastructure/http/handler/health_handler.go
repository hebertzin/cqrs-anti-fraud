package handler

import (
	"net/http"
	"time"

	"github.com/hebertzin/cqrs/internal/infrastructure/http/httpresponse"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	httpresponse.OK(w, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	httpresponse.OK(w, map[string]string{
		"status": "ready",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
