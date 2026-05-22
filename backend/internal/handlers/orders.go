package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/middlewares"
	"github.com/programvx/thecalda/backend/internal/model"
	"github.com/programvx/thecalda/backend/internal/services"
)

// OrdersHandler serves the authenticated user's order and cart endpoints.
type OrdersHandler struct {
	apiSrv    *services.ApiSrv
	ordersSrv services.OrdersSrv
}

// NewOrdersHandler constructs an OrdersHandler.
func NewOrdersHandler(apiSrv *services.ApiSrv, ordersSrv services.OrdersSrv) *OrdersHandler {
	return &OrdersHandler{apiSrv: apiSrv, ordersSrv: ordersSrv}
}

// List returns a paginated list of the authenticated user's orders. The
// optional ?type=cart|order query parameter filters by order type.
func (h *OrdersHandler) List(ctx *gin.Context) {
	authUserID, ok := middlewares.AuthUserID(ctx)
	if !ok {
		h.apiSrv.SendError(ctx, constants.ErrUnauthorized)
		return
	}
	page, pageSize := h.apiSrv.PaginationParams(ctx)

	orders, pagination, err := h.ordersSrv.List(
		ctx.Request.Context(), authUserID, ctx.Query("type"), page, pageSize,
	)
	if err != nil {
		h.apiSrv.SendError(ctx, err)
		return
	}
	h.apiSrv.SendSuccessWithMeta(ctx, orders, &model.Metadata{Pagination: pagination})
}

// Get returns one of the authenticated user's orders by uid.
func (h *OrdersHandler) Get(ctx *gin.Context) {
	authUserID, ok := middlewares.AuthUserID(ctx)
	if !ok {
		h.apiSrv.SendError(ctx, constants.ErrUnauthorized)
		return
	}

	uid, err := uuid.Parse(ctx.Param("uid"))
	if err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidUID)
		return
	}

	order, srvErr := h.ordersSrv.GetByUID(ctx.Request.Context(), authUserID, uid)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendSuccess(ctx, order)
}

// Create adds a new cart or order for the authenticated user.
func (h *OrdersHandler) Create(ctx *gin.Context) {
	authUserID, ok := middlewares.AuthUserID(ctx)
	if !ok {
		h.apiSrv.SendError(ctx, constants.ErrUnauthorized)
		return
	}

	var req model.OrderCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidBody.WithDetails(err.Error()))
		return
	}

	order, srvErr := h.ordersSrv.Create(ctx.Request.Context(), authUserID, &req)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendCreated(ctx, order)
}

// Update changes one of the authenticated user's orders.
func (h *OrdersHandler) Update(ctx *gin.Context) {
	authUserID, ok := middlewares.AuthUserID(ctx)
	if !ok {
		h.apiSrv.SendError(ctx, constants.ErrUnauthorized)
		return
	}

	uid, err := uuid.Parse(ctx.Param("uid"))
	if err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidUID)
		return
	}

	var req model.OrderUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidBody.WithDetails(err.Error()))
		return
	}

	order, srvErr := h.ordersSrv.Update(ctx.Request.Context(), authUserID, uid, &req)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendSuccess(ctx, order)
}

// Delete removes one of the authenticated user's orders.
func (h *OrdersHandler) Delete(ctx *gin.Context) {
	authUserID, ok := middlewares.AuthUserID(ctx)
	if !ok {
		h.apiSrv.SendError(ctx, constants.ErrUnauthorized)
		return
	}

	uid, err := uuid.Parse(ctx.Param("uid"))
	if err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidUID)
		return
	}

	if srvErr := h.ordersSrv.Delete(ctx.Request.Context(), authUserID, uid); srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendNoContent(ctx)
}

// Checkout places one of the authenticated user's carts as an order.
func (h *OrdersHandler) Checkout(ctx *gin.Context) {
	authUserID, ok := middlewares.AuthUserID(ctx)
	if !ok {
		h.apiSrv.SendError(ctx, constants.ErrUnauthorized)
		return
	}

	uid, err := uuid.Parse(ctx.Param("uid"))
	if err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidUID)
		return
	}

	var req model.OrderCheckout
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidBody.WithDetails(err.Error()))
		return
	}

	order, srvErr := h.ordersSrv.Checkout(ctx.Request.Context(), authUserID, uid, &req)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendSuccess(ctx, order)
}
