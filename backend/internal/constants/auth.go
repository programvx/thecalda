// Package constants holds predefined values shared across the application,
// organised per domain (one file per domain).
package constants

import (
	"net/http"

	"github.com/programvx/thecalda/backend/internal/model"
)

// Authentication and authorization errors.
var (
	ErrUnauthorized = &model.Err{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden    = &model.Err{Code: http.StatusForbidden, Message: "Forbidden"}
)
