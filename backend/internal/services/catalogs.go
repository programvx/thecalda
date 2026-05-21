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

// CatalogsSrv exposes catalog business logic.
type CatalogsSrv interface {
	Create(ctx context.Context, req *model.CatalogCreate) (*model.Catalog, *model.Err)
	List(ctx context.Context, page, pageSize int) ([]model.Catalog, *model.Pagination, *model.Err)
	GetByUID(ctx context.Context, uid uuid.UUID) (*model.Catalog, *model.Err)
	GetBySlugWithItems(ctx context.Context, slug string) (*model.CatalogWithItems, *model.Err)
	Update(ctx context.Context, uid uuid.UUID, req *model.CatalogUpdate) (*model.Catalog, *model.Err)
	Delete(ctx context.Context, uid uuid.UUID) *model.Err
}

type catalogsSrv struct {
	log      *zap.Logger
	catalogs crud.CatalogsCrud
}

// NewCatalogsSrv constructs a CatalogsSrv.
func NewCatalogsSrv(log *zap.Logger, catalogs crud.CatalogsCrud) CatalogsSrv {
	return &catalogsSrv{log: log, catalogs: catalogs}
}

// Create stores a new catalog.
func (s *catalogsSrv) Create(ctx context.Context, req *model.CatalogCreate) (*model.Catalog, *model.Err) {
	catalog := &model.Catalog{
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}
	if req.IsActive != nil {
		catalog.IsActive = *req.IsActive
	}

	if err := s.catalogs.Create(ctx, catalog); err != nil {
		return nil, s.fail("create catalog", err)
	}
	return catalog, nil
}

// List returns a page of catalogs with pagination metadata.
func (s *catalogsSrv) List(ctx context.Context, page, pageSize int) ([]model.Catalog, *model.Pagination, *model.Err) {
	catalogs, total, err := s.catalogs.List(ctx, page, pageSize)
	if err != nil {
		return nil, nil, s.fail("list catalogs", err)
	}
	return catalogs, model.NewPagination(page, pageSize, total), nil
}

// GetByUID returns a single catalog.
func (s *catalogsSrv) GetByUID(ctx context.Context, uid uuid.UUID) (*model.Catalog, *model.Err) {
	catalog, err := s.catalogs.GetByUID(ctx, uid)
	if err != nil {
		return nil, s.fail("get catalog", err)
	}
	return catalog, nil
}

// GetBySlugWithItems returns a catalog (by slug) together with its items.
func (s *catalogsSrv) GetBySlugWithItems(ctx context.Context, slug string) (*model.CatalogWithItems, *model.Err) {
	catalog, err := s.catalogs.GetBySlug(ctx, slug)
	if err != nil {
		return nil, s.fail("get catalog by slug", err)
	}

	items, err := s.catalogs.ListItems(ctx, catalog.ID)
	if err != nil {
		return nil, s.fail("list catalog items", err)
	}

	return &model.CatalogWithItems{Catalog: *catalog, Items: items}, nil
}

// Update applies the provided fields to an existing catalog.
func (s *catalogsSrv) Update(ctx context.Context, uid uuid.UUID, req *model.CatalogUpdate) (*model.Catalog, *model.Err) {
	catalog, err := s.catalogs.GetByUID(ctx, uid)
	if err != nil {
		return nil, s.fail("load catalog for update", err)
	}

	if req.Slug != nil {
		catalog.Slug = *req.Slug
	}
	if req.Name != nil {
		catalog.Name = *req.Name
	}
	if req.Description != nil {
		catalog.Description = req.Description
	}
	if req.IsActive != nil {
		catalog.IsActive = *req.IsActive
	}

	if err := s.catalogs.Update(ctx, catalog); err != nil {
		return nil, s.fail("update catalog", err)
	}
	return catalog, nil
}

// Delete removes a catalog.
func (s *catalogsSrv) Delete(ctx context.Context, uid uuid.UUID) *model.Err {
	if err := s.catalogs.Delete(ctx, uid); err != nil {
		return s.fail("delete catalog", err)
	}
	return nil
}

// fail passes domain errors through and logs/wraps everything else.
func (s *catalogsSrv) fail(op string, err error) *model.Err {
	var domainErr *model.Err
	if errors.As(err, &domainErr) {
		return domainErr
	}
	s.log.Error(op+" failed", zap.Error(err))
	return constants.ErrInternal
}
