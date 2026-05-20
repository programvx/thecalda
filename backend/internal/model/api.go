// Package model holds domain types and API DTOs shared across layers.
package model

// ApiResp is the envelope wrapping every JSON response.
type ApiResp struct {
	Data     any       `json:"data,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	Error    *ApiError `json:"error,omitempty"`
}

// ApiError is the client-facing error shape.
type ApiError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// Metadata carries response-level metadata such as pagination.
type Metadata struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination describes a paged result set.
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}
