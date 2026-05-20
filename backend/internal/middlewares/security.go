package middlewares

import "github.com/gin-gonic/gin"

// Security sets a baseline set of security response headers.
func Security() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("Referrer-Policy", "no-referrer")
		ctx.Next()
	}
}
