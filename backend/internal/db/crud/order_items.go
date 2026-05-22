package crud

import (
	"context"

	"github.com/google/uuid"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/db"
	"github.com/programvx/thecalda/backend/internal/model"
)

// OrderItemsCrud reads and writes public.order_items rows.
type OrderItemsCrud interface {
	Create(ctx context.Context, item *model.OrderItem) error
	ListByOrder(ctx context.Context, orderID int64) ([]model.OrderItem, error)
	GetByUID(ctx context.Context, userID int64, uid uuid.UUID) (*model.OrderItem, error)
	Update(ctx context.Context, item *model.OrderItem) error
	Delete(ctx context.Context, id int64) error
}

type orderItemsCrud struct {
	store *db.Store
}

// NewOrderItemsCrud constructs an OrderItemsCrud.
func NewOrderItemsCrud(store *db.Store) OrderItemsCrud {
	return &orderItemsCrud{store: store}
}

// Create inserts a new order item. The order_items recalc trigger refreshes
// the parent order's total_amount.
func (c *orderItemsCrud) Create(ctx context.Context, item *model.OrderItem) error {
	return translateError(c.store.DB.WithContext(ctx).Create(item).Error)
}

// ListByOrder returns an order's items in insertion order.
func (c *orderItemsCrud) ListByOrder(ctx context.Context, orderID int64) ([]model.OrderItem, error) {
	items := []model.OrderItem{}
	if err := c.store.DB.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("id").
		Preload("Item.Medias", preloadMedias).
		Find(&items).Error; err != nil {
		return nil, translateError(err)
	}
	return items, nil
}

// GetByUID returns one order item by uid, or constants.ErrNotFound. The join
// to orders scopes the lookup to the owning user, keeping order items private.
func (c *orderItemsCrud) GetByUID(ctx context.Context, userID int64, uid uuid.UUID) (*model.OrderItem, error) {
	var item model.OrderItem
	if err := c.store.DB.WithContext(ctx).
		Joins("JOIN orders o ON o.id = order_items.order_id").
		Where("order_items.uid = ? AND o.user_id = ?", uid, userID).
		Preload("Item.Medias", preloadMedias).
		First(&item).Error; err != nil {
		return nil, translateError(err)
	}
	return &item, nil
}

// Update persists the editable columns of an existing order item.
func (c *orderItemsCrud) Update(ctx context.Context, item *model.OrderItem) error {
	return translateError(c.store.DB.WithContext(ctx).
		Model(item).
		Select("quantity").
		Updates(item).Error)
}

// Delete removes an order item by its internal id. Callers must already have
// verified ownership (see the service layer).
func (c *orderItemsCrud) Delete(ctx context.Context, id int64) error {
	res := c.store.DB.WithContext(ctx).Delete(&model.OrderItem{}, id)
	if res.Error != nil {
		return translateError(res.Error)
	}
	if res.RowsAffected == 0 {
		return constants.ErrNotFound
	}
	return nil
}
