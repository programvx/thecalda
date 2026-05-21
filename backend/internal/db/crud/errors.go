package crud

import (
	"errors"

	"gorm.io/gorm"

	"github.com/programvx/thecalda/backend/internal/constants"
)

// translateError maps GORM/driver errors to domain errors. Errors it does not
// recognise are returned unchanged for the service layer to log and wrap.
func translateError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return constants.ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return constants.ErrConflict
	default:
		return err
	}
}
