package crud

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/db"
	"github.com/programvx/thecalda/backend/internal/model"
)

// OrdersCrud reads and writes public.orders rows and their order_items.
type OrdersCrud interface {
	Create(ctx context.Context, order *model.Order) error
	List(ctx context.Context, userID int64, orderType string, page, pageSize int) ([]model.Order, int64, error)
	GetByUID(ctx context.Context, userID int64, uid uuid.UUID) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, userID int64, uid uuid.UUID) error
	ItemsByUIDs(ctx context.Context, uids []uuid.UUID) ([]model.Item, error)
}

type ordersCrud struct {
	store *db.Store
}

// NewOrdersCrud constructs an OrdersCrud.
func NewOrdersCrud(store *db.Store) OrdersCrud {
	return &ordersCrud{store: store}
}

// preloadItems loads an order's items in insertion order.
func preloadItems(tx *gorm.DB) *gorm.DB {
	return tx.Order("order_items.id")
}

// Create inserts an order together with its order_items in one transaction.
func (c *ordersCrud) Create(ctx context.Context, order *model.Order) error {
	return translateError(c.store.DB.WithContext(ctx).Create(order).Error)
}

// List returns a page of a user's orders (newest first) with their items, plus
// the total row count. orderType, when non-empty, filters by 'cart' or 'order'.
func (c *ordersCrud) List(ctx context.Context, userID int64, orderType string, page, pageSize int) ([]model.Order, int64, error) {
	orders := []model.Order{}
	var total int64

	countQ := c.store.DB.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)
	listQ := c.store.DB.WithContext(ctx).Where("user_id = ?", userID)
	if orderType != "" {
		countQ = countQ.Where("type = ?", orderType)
		listQ = listQ.Where("type = ?", orderType)
	}

	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, translateError(err)
	}
	if err := listQ.
		Order("id DESC").
		Limit(pageSize).
		Offset((page-1)*pageSize).
		Preload("Items", preloadItems).
		Find(&orders).Error; err != nil {
		return nil, 0, translateError(err)
	}
	return orders, total, nil
}

// GetByUID returns one of the user's orders with its items, or
// constants.ErrNotFound. Scoping by user_id keeps orders private per owner.
func (c *ordersCrud) GetByUID(ctx context.Context, userID int64, uid uuid.UUID) (*model.Order, error) {
	var order model.Order
	if err := c.store.DB.WithContext(ctx).
		Where("uid = ? AND user_id = ?", uid, userID).
		Preload("Items", preloadItems).
		First(&order).Error; err != nil {
		return nil, translateError(err)
	}
	return &order, nil
}

// Update persists the editable columns of an existing order.
func (c *ordersCrud) Update(ctx context.Context, order *model.Order) error {
	return translateError(c.store.DB.WithContext(ctx).
		Model(order).
		Select("status", "notes").
		Updates(order).Error)
}

// Delete removes one of the user's orders (order_items cascade), or returns
// constants.ErrNotFound when no such order exists for that user.
func (c *ordersCrud) Delete(ctx context.Context, userID int64, uid uuid.UUID) error {
	res := c.store.DB.WithContext(ctx).
		Where("uid = ? AND user_id = ?", uid, userID).
		Delete(&model.Order{})
	if res.Error != nil {
		return translateError(res.Error)
	}
	if res.RowsAffected == 0 {
		return constants.ErrNotFound
	}
	return nil
}

// ItemsByUIDs loads the active catalog items for the given uids. The orders
// service uses it to snapshot item details onto order_items at create time.
func (c *ordersCrud) ItemsByUIDs(ctx context.Context, uids []uuid.UUID) ([]model.Item, error) {
	items := []model.Item{}
	if len(uids) == 0 {
		return items, nil
	}
	if err := c.store.DB.WithContext(ctx).
		Where("uid IN ? AND is_active", uids).
		Find(&items).Error; err != nil {
		return nil, translateError(err)
	}
	return items, nil
}
