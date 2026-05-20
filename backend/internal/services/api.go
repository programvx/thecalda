// Package services holds business logic and shared HTTP response helpers.
package services

import (
	"net/http"

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
