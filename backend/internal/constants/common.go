package constants

import (
	"net/http"

	"github.com/programvx/thecalda/backend/internal/model"
)

// Common errors shared across domains.
var (
	ErrNotFound    = &model.Err{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrInvalidUID  = &model.Err{Code: http.StatusBadRequest, Message: "Invalid resource id"}
	ErrInvalidBody = &model.Err{Code: http.StatusBadRequest, Message: "Invalid request body"}
	ErrConflict    = &model.Err{Code: http.StatusConflict, Message: "Resource already exists"}
	ErrInternal    = &model.Err{Code: http.StatusInternalServerError, Message: "Internal server error"}
)
