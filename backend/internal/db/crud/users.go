// Package crud holds the database repository implementations (GORM-backed).
package crud

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/db"
	"github.com/programvx/thecalda/backend/internal/model"
)

// UsersCrud reads and writes public.users rows.
type UsersCrud interface {
	GetByAuthUserID(ctx context.Context, authUserID uuid.UUID) (*model.User, error)
}

type usersCrud struct {
	store *db.Store
}

// NewUsersCrud constructs a UsersCrud.
func NewUsersCrud(store *db.Store) UsersCrud {
	return &usersCrud{store: store}
}

// GetByAuthUserID returns the profile linked to the given Supabase Auth user,
// or constants.ErrNotFound when no such profile exists.
func (c *usersCrud) GetByAuthUserID(ctx context.Context, authUserID uuid.UUID) (*model.User, error) {
	var user model.User

	err := c.store.DB.WithContext(ctx).
		Where("auth_user_id = ?", authUserID).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
