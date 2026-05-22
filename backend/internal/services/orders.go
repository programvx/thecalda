package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/db/crud"
	"github.com/programvx/thecalda/backend/internal/model"
)

// Order type and status values, mirroring the order_type / order_status enums
// in migration 003_orders.
const (
	orderTypeCart  = "cart"
	orderTypeOrder = "order"

	orderStatusNotApplicable = "not_applicable"
	orderStatusPending       = "pending"
)

// settableOrderStatuses are the lifecycle statuses a client may assign to an
// order via Update ('not_applicable' is internal to carts).
var settableOrderStatuses = map[string]bool{
	"pending":   true,
	"paid":      true,
	"shipped":   true,
	"delivered": true,
	"cancelled": true,
}

// OrdersSrv exposes order and cart business logic, scoped to the
// authenticated user.
type OrdersSrv interface {
	Create(ctx context.Context, authUserID uuid.UUID, req *model.OrderCreate) (*model.Order, *model.Err)
	List(ctx context.Context, authUserID uuid.UUID, orderType string, page, pageSize int) ([]model.Order, *model.Pagination, *model.Err)
	GetByUID(ctx context.Context, authUserID, uid uuid.UUID) (*model.Order, *model.Err)
	Update(ctx context.Context, authUserID, uid uuid.UUID, req *model.OrderUpdate) (*model.Order, *model.Err)
	Delete(ctx context.Context, authUserID, uid uuid.UUID) *model.Err
}

type ordersSrv struct {
	log    *zap.Logger
	orders crud.OrdersCrud
	users  crud.UsersCrud
}

// NewOrdersSrv constructs an OrdersSrv. It needs UsersCrud to resolve the
// authenticated Supabase user to the public.users row that owns the orders.
func NewOrdersSrv(log *zap.Logger, orders crud.OrdersCrud, users crud.UsersCrud) OrdersSrv {
	return &ordersSrv{log: log, orders: orders, users: users}
}

// Create stores a new cart or order for the user, snapshotting catalog items.
func (s *ordersSrv) Create(ctx context.Context, authUserID uuid.UUID, req *model.OrderCreate) (*model.Order, *model.Err) {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return nil, srvErr
	}

	orderType := orderTypeCart
	if req.Type != "" {
		if req.Type != orderTypeCart && req.Type != orderTypeOrder {
			return nil, constants.ErrInvalidBody.WithDetails("type must be 'cart' or 'order'")
		}
		orderType = req.Type
	}

	status := orderStatusNotApplicable
	if orderType == orderTypeOrder {
		status = orderStatusPending
	}

	// Collect requested quantities, rejecting an item listed more than once.
	qtyByUID := make(map[uuid.UUID]int, len(req.Items))
	uids := make([]uuid.UUID, 0, len(req.Items))
	for _, line := range req.Items {
		if _, dup := qtyByUID[line.ItemUID]; dup {
			return nil, constants.ErrInvalidBody.WithDetails("an item is listed more than once")
		}
		qtyByUID[line.ItemUID] = line.Quantity
		uids = append(uids, line.ItemUID)
	}

	catalogItems, err := s.orders.ItemsByUIDs(ctx, uids)
	if err != nil {
		return nil, s.fail("load order items", err)
	}
	if len(catalogItems) != len(uids) {
		return nil, constants.ErrInvalidBody.WithDetails("one or more items do not exist or are unavailable")
	}

	// Snapshot each catalog item onto an order line.
	orderItems := make([]model.OrderItem, 0, len(catalogItems))
	for _, it := range catalogItems {
		orderItems = append(orderItems, model.OrderItem{
			ItemID:       it.ID,
			ItemName:     it.Name,
			Quantity:     qtyByUID[it.UID],
			UnitPrice:    it.Price,
			UnitDiscount: it.Discount,
		})
	}

	order := &model.Order{
		UserID:     userID,
		Type:       orderType,
		Status:     status,
		CurrencyID: catalogItems[0].CurrencyID,
		Notes:      req.Notes,
		Items:      orderItems,
	}
	if err := s.orders.Create(ctx, order); err != nil {
		return nil, s.fail("create order", err)
	}

	// Re-fetch so the response carries the database-computed totals.
	created, err := s.orders.GetByUID(ctx, userID, order.UID)
	if err != nil {
		return nil, s.fail("reload created order", err)
	}
	return created, nil
}

// List returns a page of the user's orders with pagination metadata.
func (s *ordersSrv) List(ctx context.Context, authUserID uuid.UUID, orderType string, page, pageSize int) ([]model.Order, *model.Pagination, *model.Err) {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return nil, nil, srvErr
	}

	// Only the known order types filter the list; anything else lists all.
	filter := ""
	if orderType == orderTypeCart || orderType == orderTypeOrder {
		filter = orderType
	}

	orders, total, err := s.orders.List(ctx, userID, filter, page, pageSize)
	if err != nil {
		return nil, nil, s.fail("list orders", err)
	}
	return orders, model.NewPagination(page, pageSize, total), nil
}

// GetByUID returns one of the user's orders.
func (s *ordersSrv) GetByUID(ctx context.Context, authUserID, uid uuid.UUID) (*model.Order, *model.Err) {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return nil, srvErr
	}

	order, err := s.orders.GetByUID(ctx, userID, uid)
	if err != nil {
		return nil, s.fail("get order", err)
	}
	return order, nil
}

// Update applies the provided fields to one of the user's orders.
func (s *ordersSrv) Update(ctx context.Context, authUserID, uid uuid.UUID, req *model.OrderUpdate) (*model.Order, *model.Err) {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return nil, srvErr
	}

	order, err := s.orders.GetByUID(ctx, userID, uid)
	if err != nil {
		return nil, s.fail("load order for update", err)
	}

	if req.Status != nil {
		if order.Type == orderTypeCart {
			return nil, constants.ErrInvalidBody.WithDetails("a cart has no settable status")
		}
		if !settableOrderStatuses[*req.Status] {
			return nil, constants.ErrInvalidBody.WithDetails("invalid order status")
		}
		order.Status = *req.Status
	}
	if req.Notes != nil {
		order.Notes = req.Notes
	}

	if err := s.orders.Update(ctx, order); err != nil {
		return nil, s.fail("update order", err)
	}
	return order, nil
}

// Delete removes one of the user's orders.
func (s *ordersSrv) Delete(ctx context.Context, authUserID, uid uuid.UUID) *model.Err {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return srvErr
	}
	if err := s.orders.Delete(ctx, userID, uid); err != nil {
		return s.fail("delete order", err)
	}
	return nil
}

// resolveUserID maps the authenticated Supabase user to the public.users id
// that owns their orders.
func (s *ordersSrv) resolveUserID(ctx context.Context, authUserID uuid.UUID) (int64, *model.Err) {
	user, err := s.users.GetByAuthUserID(ctx, authUserID)
	if err != nil {
		return 0, s.fail("resolve user", err)
	}
	return user.ID, nil
}

// fail passes domain errors through and logs/wraps everything else.
func (s *ordersSrv) fail(op string, err error) *model.Err {
	var domainErr *model.Err
	if errors.As(err, &domainErr) {
		return domainErr
	}
	s.log.Error(op+" failed", zap.Error(err))
	return constants.ErrInternal
}
