package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db        *pgxpool.Pool
	startTime time.Time
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
	}
}

type HealthResponse struct {
	Status   string `json:"status"`
	Uptime   string `json:"uptime"`
	Database string `json:"database"`
	Service  string `json:"service"`
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	dbStatus := "connected"

	if h.db != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := h.db.Ping(ctx); err != nil {
			dbStatus = "disconnected: " + err.Error()
		}
	} else {
		dbStatus = "not_initialized"
	}

	status := "healthy"
	if dbStatus != "connected" {
		status = "degraded"
	}

	JSON(w, http.StatusOK, HealthResponse{
		Status:   status,
		Uptime:   time.Since(h.startTime).String(),
		Database: dbStatus,
		Service:  "semard-control-center-api",
	})
}
