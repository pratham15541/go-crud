package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pratham15541/go-crud/internal/models"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db        *sql.DB
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
	}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check database connection
	dbStatus := "healthy"
	dbError := ""
	if err := h.db.Ping(); err != nil {
		dbStatus = "unhealthy"
		dbError = err.Error()
	}

	// Calculate uptime
	uptime := time.Since(h.startTime)

	// Create health response
	healthResp := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
		Checks: map[string]interface{}{
			"database": map[string]interface{}{
				"status": dbStatus,
				"error":  dbError,
			},
			"memory": map[string]interface{}{
				"status": "healthy",
			},
		},
	}

	// If database is unhealthy, mark overall status as unhealthy
	if dbStatus == "unhealthy" {
		healthResp.Status = "unhealthy"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(healthResp)
}