package constants

import (
	"net/http"

	"github.com/programvx/thecalda/backend/internal/model"
)

// Common errors shared across domains.
var (
	ErrNotFound    = &model.Err{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrInvalidBody = &model.Err{Code: http.StatusBadRequest, Message: "Invalid request body"}
	ErrInternal    = &model.Err{Code: http.StatusInternalServerError, Message: "Internal server error"}
)
