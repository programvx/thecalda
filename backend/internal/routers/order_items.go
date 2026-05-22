package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/programvx/thecalda/backend/internal/handlers"
)

// NewOrderItemsRouter registers the order-item routes. Order items are private
// to the authenticated user, so the given router must already enforce
// authentication (see the /api group in cmd/api).
func NewOrderItemsRouter(r gin.IRouter, h *handlers.OrderItemsHandler) {
	items := r.Group("/order-items")
	items.GET("", h.List)
	items.GET("/:uid", h.Get)
	items.POST("", h.Create)
	items.PUT("/:uid", h.Update)
	items.DELETE("/:uid", h.Delete)
}
