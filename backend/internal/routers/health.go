// Package routers registers HTTP handlers onto the Gin router. The Gin engine
// and global middleware are constructed in cmd/api; this package only maps
// routes to handlers, split per domain.
package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/programvx/thecalda/backend/internal/handlers"
)

// NewHealthRouter registers the health-check route.
func NewHealthRouter(r gin.IRouter, h *handlers.HealthHandler) {
	r.GET("/health", h.Check)
}
