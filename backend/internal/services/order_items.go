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

// OrderItemsSrv exposes order-item business logic, scoped to the
// authenticated user.
type OrderItemsSrv interface {
	Create(ctx context.Context, authUserID uuid.UUID, req *model.OrderItemCreate) (*model.OrderItem, *model.Err)
	ListByOrder(ctx context.Context, authUserID, orderUID uuid.UUID) ([]model.OrderItem, *model.Err)
	GetByUID(ctx context.Context, authUserID, uid uuid.UUID) (*model.OrderItem, *model.Err)
	Update(ctx context.Context, authUserID, uid uuid.UUID, req *model.OrderItemUpdate) (*model.OrderItem, *model.Err)
	Delete(ctx context.Context, authUserID, uid uuid.UUID) *model.Err
}

type orderItemsSrv struct {
	log        *zap.Logger
	orderItems crud.OrderItemsCrud
	orders     crud.OrdersCrud
	users      crud.UsersCrud
}

// NewOrderItemsSrv constructs an OrderItemsSrv. It needs OrdersCrud to verify
// the parent order belongs to the caller and to snapshot catalog items, and
// UsersCrud to resolve the authenticated Supabase user.
func NewOrderItemsSrv(log *zap.Logger, orderItems crud.OrderItemsCrud, orders crud.OrdersCrud, users crud.UsersCrud) OrderItemsSrv {
	return &orderItemsSrv{log: log, orderItems: orderItems, orders: orders, users: users}
}

// Create adds a catalog item to one of the user's orders, snapshotting its
// name and price onto the new line.
func (s *orderItemsSrv) Create(ctx context.Context, authUserID uuid.UUID, req *model.OrderItemCreate) (*model.OrderItem, *model.Err) {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return nil, srvErr
	}

	if req.OrderUID == uuid.Nil || req.ItemUID == uuid.Nil {
		return nil, constants.ErrInvalidBody.WithDetails("orderUid and itemUid are required")
	}

	// The order must exist and belong to the caller.
	order, err := s.orders.GetByUID(ctx, userID, req.OrderUID)
	if err != nil {
		return nil, s.fail("load order", err)
	}

	// Snapshot the catalog item.
	items, err := s.orders.ItemsByUIDs(ctx, []uuid.UUID{req.ItemUID})
	if err != nil {
		return nil, s.fail("load catalog item", err)
	}
	if len(items) != 1 {
		return nil, constants.ErrInvalidBody.WithDetails("item does not exist or is unavailable")
	}
	catalogItem := items[0]

	orderItem := &model.OrderItem{
		OrderID:      order.ID,
		ItemID:       catalogItem.ID,
		ItemName:     catalogItem.Name,
		Quantity:     req.Quantity,
		UnitPrice:    catalogItem.Price,
		UnitDiscount: catalogItem.Discount,
	}
	if err := s.orderItems.Create(ctx, orderItem); err != nil {
		return nil, s.fail("create order item", err)
	}

	// Re-fetch so the response carries the database-computed line totals.
	created, err := s.orderItems.GetByUID(ctx, userID, orderItem.UID)
	if err != nil {
		return nil, s.fail("reload created order item", err)
	}
	return created, nil
}

// ListByOrder returns the items of one of the user's orders.
func (s *orderItemsSrv) ListByOrder(ctx context.Context, authUserID, orderUID uuid.UUID) ([]model.OrderItem, *model.Err) {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return nil, srvErr
	}

	if orderUID == uuid.Nil {
		return nil, constants.ErrInvalidBody.WithDetails("orderUid is required")
	}

	order, err := s.orders.GetByUID(ctx, userID, orderUID)
	if err != nil {
		return nil, s.fail("load order", err)
	}

	items, err := s.orderItems.ListByOrder(ctx, order.ID)
	if err != nil {
		return nil, s.fail("list order items", err)
	}
	return items, nil
}

// GetByUID returns one of the user's order items.
func (s *orderItemsSrv) GetByUID(ctx context.Context, authUserID, uid uuid.UUID) (*model.OrderItem, *model.Err) {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return nil, srvErr
	}

	item, err := s.orderItems.GetByUID(ctx, userID, uid)
	if err != nil {
		return nil, s.fail("get order item", err)
	}
	return item, nil
}

// Update changes the quantity of one of the user's order items.
func (s *orderItemsSrv) Update(ctx context.Context, authUserID, uid uuid.UUID, req *model.OrderItemUpdate) (*model.OrderItem, *model.Err) {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return nil, srvErr
	}

	item, err := s.orderItems.GetByUID(ctx, userID, uid)
	if err != nil {
		return nil, s.fail("load order item for update", err)
	}

	item.Quantity = req.Quantity
	if err := s.orderItems.Update(ctx, item); err != nil {
		return nil, s.fail("update order item", err)
	}

	// Re-fetch so the response carries the recomputed line total.
	updated, err := s.orderItems.GetByUID(ctx, userID, uid)
	if err != nil {
		return nil, s.fail("reload updated order item", err)
	}
	return updated, nil
}

// Delete removes one of the user's order items.
func (s *orderItemsSrv) Delete(ctx context.Context, authUserID, uid uuid.UUID) *model.Err {
	userID, srvErr := s.resolveUserID(ctx, authUserID)
	if srvErr != nil {
		return srvErr
	}

	item, err := s.orderItems.GetByUID(ctx, userID, uid)
	if err != nil {
		return s.fail("load order item for delete", err)
	}
	if err := s.orderItems.Delete(ctx, item.ID); err != nil {
		return s.fail("delete order item", err)
	}
	return nil
}

// resolveUserID maps the authenticated Supabase user to the public.users id
// that owns their orders.
func (s *orderItemsSrv) resolveUserID(ctx context.Context, authUserID uuid.UUID) (int64, *model.Err) {
	user, err := s.users.GetByAuthUserID(ctx, authUserID)
	if err != nil {
		return 0, s.fail("resolve user", err)
	}
	return user.ID, nil
}

// fail passes domain errors through and logs/wraps everything else.
func (s *orderItemsSrv) fail(op string, err error) *model.Err {
	var domainErr *model.Err
	if errors.As(err, &domainErr) {
		return domainErr
	}
	s.log.Error(op+" failed", zap.Error(err))
	return constants.ErrInternal
}
