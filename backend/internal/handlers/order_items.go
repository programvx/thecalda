package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/middlewares"
	"github.com/programvx/thecalda/backend/internal/model"
	"github.com/programvx/thecalda/backend/internal/services"
)

// OrderItemsHandler serves the authenticated user's order-item endpoints.
type OrderItemsHandler struct {
	apiSrv        *services.ApiSrv
	orderItemsSrv services.OrderItemsSrv
}

// NewOrderItemsHandler constructs an OrderItemsHandler.
func NewOrderItemsHandler(apiSrv *services.ApiSrv, orderItemsSrv services.OrderItemsSrv) *OrderItemsHandler {
	return &OrderItemsHandler{apiSrv: apiSrv, orderItemsSrv: orderItemsSrv}
}

// List returns the items of one of the authenticated user's orders. The
// orderUid query parameter is required.
func (h *OrderItemsHandler) List(ctx *gin.Context) {
	authUserID, ok := middlewares.AuthUserID(ctx)
	if !ok {
		h.apiSrv.SendError(ctx, constants.ErrUnauthorized)
		return
	}

	orderUID, err := uuid.Parse(ctx.Query("orderUid"))
	if err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidBody.WithDetails("orderUid query parameter is required"))
		return
	}

	items, srvErr := h.orderItemsSrv.ListByOrder(ctx.Request.Context(), authUserID, orderUID)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendSuccess(ctx, items)
}

// Get returns one of the authenticated user's order items by uid.
func (h *OrderItemsHandler) Get(ctx *gin.Context) {
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

	item, srvErr := h.orderItemsSrv.GetByUID(ctx.Request.Context(), authUserID, uid)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendSuccess(ctx, item)
}

// Create adds a catalog item to one of the authenticated user's orders.
func (h *OrderItemsHandler) Create(ctx *gin.Context) {
	authUserID, ok := middlewares.AuthUserID(ctx)
	if !ok {
		h.apiSrv.SendError(ctx, constants.ErrUnauthorized)
		return
	}

	var req model.OrderItemCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidBody.WithDetails(err.Error()))
		return
	}

	item, srvErr := h.orderItemsSrv.Create(ctx.Request.Context(), authUserID, &req)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendCreated(ctx, item)
}

// Update changes the quantity of one of the authenticated user's order items.
func (h *OrderItemsHandler) Update(ctx *gin.Context) {
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

	var req model.OrderItemUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidBody.WithDetails(err.Error()))
		return
	}

	item, srvErr := h.orderItemsSrv.Update(ctx.Request.Context(), authUserID, uid, &req)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendSuccess(ctx, item)
}

// Delete removes one of the authenticated user's order items.
func (h *OrderItemsHandler) Delete(ctx *gin.Context) {
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

	if srvErr := h.orderItemsSrv.Delete(ctx.Request.Context(), authUserID, uid); srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendNoContent(ctx)
}
