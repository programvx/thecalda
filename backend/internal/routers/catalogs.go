package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/programvx/thecalda/backend/internal/handlers"
)

// NewCatalogsPublicRouter registers the public catalog read routes.
func NewCatalogsPublicRouter(r gin.IRouter, h *handlers.CatalogsHandler) {
	catalogs := r.Group("/catalogs")
	catalogs.GET("", h.List)
	catalogs.GET("/:slug", h.GetBySlug)
}

// NewCatalogsPrivateRouter registers the authenticated catalog write routes.
func NewCatalogsPrivateRouter(r gin.IRouter, h *handlers.CatalogsHandler) {
	catalogs := r.Group("/catalogs")
	catalogs.POST("", h.Create)
	catalogs.PUT("/:uid", h.Update)
	catalogs.DELETE("/:uid", h.Delete)
}
