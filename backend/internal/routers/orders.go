package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/programvx/thecalda/backend/internal/handlers"
)

// NewOrdersRouter registers the order routes. Orders are private to the
// authenticated user, so the given router must already enforce authentication
// (see the /api group in cmd/api).
func NewOrdersRouter(r gin.IRouter, h *handlers.OrdersHandler) {
	orders := r.Group("/orders")
	orders.GET("", h.List)
	orders.GET("/:uid", h.Get)
	orders.POST("", h.Create)
	orders.POST("/:uid/checkout", h.Checkout)
	orders.PUT("/:uid", h.Update)
	orders.DELETE("/:uid", h.Delete)
}
