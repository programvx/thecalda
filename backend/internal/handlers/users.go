package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/middlewares"
	"github.com/programvx/thecalda/backend/internal/services"
)

// UsersHandler serves user endpoints.
type UsersHandler struct {
	apiSrv   *services.ApiSrv
	usersSrv services.UsersSrv
}

// NewUsersHandler constructs a UsersHandler.
func NewUsersHandler(apiSrv *services.ApiSrv, usersSrv services.UsersSrv) *UsersHandler {
	return &UsersHandler{apiSrv: apiSrv, usersSrv: usersSrv}
}

// Me returns the authenticated caller's profile.
func (h *UsersHandler) Me(ctx *gin.Context) {
	authUserID, ok := middlewares.AuthUserID(ctx)
	if !ok {
		h.apiSrv.SendError(ctx, constants.ErrUnauthorized)
		return
	}

	user, err := h.usersSrv.GetByAuthUserID(ctx.Request.Context(), authUserID)
	if err != nil {
		h.apiSrv.SendError(ctx, err)
		return
	}

	h.apiSrv.SendSuccess(ctx, user)
}
