package crud

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/db"
	"github.com/programvx/thecalda/backend/internal/model"
)

// CatalogsCrud reads and writes public.catalogs rows.
type CatalogsCrud interface {
	Create(ctx context.Context, catalog *model.Catalog) error
	List(ctx context.Context, page, pageSize int) ([]model.Catalog, int64, error)
	GetByUID(ctx context.Context, uid uuid.UUID) (*model.Catalog, error)
	GetBySlug(ctx context.Context, slug string) (*model.Catalog, error)
	Update(ctx context.Context, catalog *model.Catalog) error
	Delete(ctx context.Context, uid uuid.UUID) error
	ListItems(ctx context.Context, catalogID int64) ([]model.Item, error)
}

type catalogsCrud struct {
	store *db.Store
}

// NewCatalogsCrud constructs a CatalogsCrud.
func NewCatalogsCrud(store *db.Store) CatalogsCrud {
	return &catalogsCrud{store: store}
}

// Create inserts a new catalog, populating its generated fields.
func (c *catalogsCrud) Create(ctx context.Context, catalog *model.Catalog) error {
	return translateError(c.store.DB.WithContext(ctx).Create(catalog).Error)
}

// List returns a page of catalogs ordered by id, plus the total row count.
func (c *catalogsCrud) List(ctx context.Context, page, pageSize int) ([]model.Catalog, int64, error) {
	catalogs := []model.Catalog{}
	var total int64

	q := c.store.DB.WithContext(ctx)
	if err := q.Model(&model.Catalog{}).Count(&total).Error; err != nil {
		return nil, 0, translateError(err)
	}
	if err := q.Order("id").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&catalogs).Error; err != nil {
		return nil, 0, translateError(err)
	}
	return catalogs, total, nil
}

// GetByUID returns the catalog with the given uid, or constants.ErrNotFound.
func (c *catalogsCrud) GetByUID(ctx context.Context, uid uuid.UUID) (*model.Catalog, error) {
	var catalog model.Catalog
	if err := c.store.DB.WithContext(ctx).
		Where("uid = ?", uid).
		First(&catalog).Error; err != nil {
		return nil, translateError(err)
	}
	return &catalog, nil
}

// GetBySlug returns the catalog with the given slug, or constants.ErrNotFound.
func (c *catalogsCrud) GetBySlug(ctx context.Context, slug string) (*model.Catalog, error) {
	var catalog model.Catalog
	if err := c.store.DB.WithContext(ctx).
		Where("slug = ?", slug).
		First(&catalog).Error; err != nil {
		return nil, translateError(err)
	}
	return &catalog, nil
}

// Update persists changes to an existing catalog.
func (c *catalogsCrud) Update(ctx context.Context, catalog *model.Catalog) error {
	return translateError(c.store.DB.WithContext(ctx).Save(catalog).Error)
}

// Delete removes the catalog with the given uid, or returns constants.ErrNotFound.
func (c *catalogsCrud) Delete(ctx context.Context, uid uuid.UUID) error {
	res := c.store.DB.WithContext(ctx).
		Where("uid = ?", uid).
		Delete(&model.Catalog{})
	if res.Error != nil {
		return translateError(res.Error)
	}
	if res.RowsAffected == 0 {
		return constants.ErrNotFound
	}
	return nil
}

// ListItems returns the active items belonging to a catalog, ordered by their
// position within that catalog. Each item's stock, media (ordered so the
// primary, lowest-position media comes first), and properties are preloaded.
func (c *catalogsCrud) ListItems(ctx context.Context, catalogID int64) ([]model.Item, error) {
	items := []model.Item{}
	err := c.store.DB.WithContext(ctx).
		Joins("JOIN catalog_items ci ON ci.item_id = items.id").
		Where("ci.catalog_id = ? AND items.is_active", catalogID).
		Order("ci.position").
		Preload("Stock").
		Preload("Medias", func(db *gorm.DB) *gorm.DB {
			return db.Order("item_medias.position")
		}).
		Preload("Properties", func(db *gorm.DB) *gorm.DB {
			return db.Order("item_properties.id")
		}).
		Find(&items).Error
	if err != nil {
		return nil, translateError(err)
	}
	return items, nil
}
