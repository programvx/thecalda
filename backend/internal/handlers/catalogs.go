package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/model"
	"github.com/programvx/thecalda/backend/internal/services"
)

// CatalogsHandler serves catalog CRUD endpoints.
type CatalogsHandler struct {
	apiSrv      *services.ApiSrv
	catalogsSrv services.CatalogsSrv
}

// NewCatalogsHandler constructs a CatalogsHandler.
func NewCatalogsHandler(apiSrv *services.ApiSrv, catalogsSrv services.CatalogsSrv) *CatalogsHandler {
	return &CatalogsHandler{apiSrv: apiSrv, catalogsSrv: catalogsSrv}
}

// List returns a paginated list of catalogs.
func (h *CatalogsHandler) List(ctx *gin.Context) {
	page, pageSize := h.apiSrv.PaginationParams(ctx)

	catalogs, pagination, err := h.catalogsSrv.List(ctx.Request.Context(), page, pageSize)
	if err != nil {
		h.apiSrv.SendError(ctx, err)
		return
	}
	h.apiSrv.SendSuccessWithMeta(ctx, catalogs, &model.Metadata{Pagination: pagination})
}

// GetBySlug returns a single catalog (by slug) together with its items.
func (h *CatalogsHandler) GetBySlug(ctx *gin.Context) {
	catalog, srvErr := h.catalogsSrv.GetBySlugWithItems(ctx.Request.Context(), ctx.Param("slug"))
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendSuccess(ctx, catalog)
}

// Create adds a new catalog.
func (h *CatalogsHandler) Create(ctx *gin.Context) {
	var req model.CatalogCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidBody.WithDetails(err.Error()))
		return
	}

	catalog, srvErr := h.catalogsSrv.Create(ctx.Request.Context(), &req)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendCreated(ctx, catalog)
}

// Update changes an existing catalog.
func (h *CatalogsHandler) Update(ctx *gin.Context) {
	uid, err := uuid.Parse(ctx.Param("uid"))
	if err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidUID)
		return
	}

	var req model.CatalogUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidBody.WithDetails(err.Error()))
		return
	}

	catalog, srvErr := h.catalogsSrv.Update(ctx.Request.Context(), uid, &req)
	if srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendSuccess(ctx, catalog)
}

// Delete removes a catalog by uid.
func (h *CatalogsHandler) Delete(ctx *gin.Context) {
	uid, err := uuid.Parse(ctx.Param("uid"))
	if err != nil {
		h.apiSrv.SendError(ctx, constants.ErrInvalidUID)
		return
	}

	if srvErr := h.catalogsSrv.Delete(ctx.Request.Context(), uid); srvErr != nil {
		h.apiSrv.SendError(ctx, srvErr)
		return
	}
	h.apiSrv.SendNoContent(ctx)
}
