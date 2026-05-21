// Package services holds business logic and shared HTTP response helpers.
package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/programvx/thecalda/backend/internal/model"
)

// ApiSrv centralises JSON response formatting so handlers stay consistent.
type ApiSrv struct{}

// NewApiSrv constructs an ApiSrv.
func NewApiSrv() *ApiSrv {
	return &ApiSrv{}
}

// SendSuccess writes a 200 response wrapping data.
func (s *ApiSrv) SendSuccess(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, model.ApiResp{Data: data})
}

// SendSuccessWithMeta writes a 200 response with response-level metadata.
func (s *ApiSrv) SendSuccessWithMeta(ctx *gin.Context, data any, meta *model.Metadata) {
	ctx.JSON(http.StatusOK, model.ApiResp{Data: data, Metadata: meta})
}

// SendCreated writes a 201 response wrapping data.
func (s *ApiSrv) SendCreated(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusCreated, model.ApiResp{Data: data})
}

// SendNoContent writes a 204 response.
func (s *ApiSrv) SendNoContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

// SendError writes an error response derived from a domain error.
func (s *ApiSrv) SendError(ctx *gin.Context, err *model.Err) {
	ctx.AbortWithStatusJSON(err.Code, model.ApiResp{
		Error: &model.ApiError{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		},
	})
}

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// PaginationParams reads `page` and `pageSize` query parameters, applying
// defaults and clamping them to sane bounds.
func (s *ApiSrv) PaginationParams(ctx *gin.Context) (page, pageSize int) {
	page = queryInt(ctx, "page", defaultPage)
	if page < 1 {
		page = defaultPage
	}

	pageSize = queryInt(ctx, "pageSize", defaultPageSize)
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

// queryInt reads an integer query parameter, falling back to def when the
// parameter is absent or not a valid integer.
func queryInt(ctx *gin.Context, key string, def int) int {
	raw := ctx.Query(key)
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return n
}
