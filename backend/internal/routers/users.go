package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/programvx/thecalda/backend/internal/handlers"
)

// NewUsersRouter registers the user routes. The given router is expected
// to already enforce authentication (see the /api group in cmd/api).
func NewUsersRouter(r gin.IRouter, h *handlers.UsersHandler) {
	r.GET("/me", h.Me)
}
