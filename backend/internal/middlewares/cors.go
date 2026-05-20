// Package middlewares holds Gin middleware.
package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/programvx/thecalda/backend/internal/settings"
)

// CORS returns a middleware that applies the configured CORS policy.
func CORS(cfg *settings.Settings) gin.HandlerFunc {
	allowed := cfg.CORSAllowOrigins

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if origin != "" && slices.Contains(allowed, origin) {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			ctx.Header("Vary", "Origin")
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
