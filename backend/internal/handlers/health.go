// Package handlers holds HTTP request handlers.
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/programvx/thecalda/backend/internal/db"
)

// HealthHandler reports service liveness.
type HealthHandler struct {
	store *db.Store
}

// NewHealthHandler constructs a HealthHandler.
func NewHealthHandler(store *db.Store) *HealthHandler {
	return &HealthHandler{store: store}
}

// Check reports whether the API and its database are reachable.
func (h *HealthHandler) Check(ctx *gin.Context) {
	dbStatus := "ok"
	httpStatus := http.StatusOK

	pingCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()
	if err := h.store.Ping(pingCtx); err != nil {
		dbStatus = "unreachable"
		httpStatus = http.StatusServiceUnavailable
	}

	ctx.JSON(httpStatus, gin.H{
		"status":   "ok",
		"database": dbStatus,
	})
}
