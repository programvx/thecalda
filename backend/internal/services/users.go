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

// UsersSrv exposes user-related business logic.
type UsersSrv interface {
	GetByAuthUserID(ctx context.Context, authUserID uuid.UUID) (*model.User, *model.Err)
}

type usersSrv struct {
	log   *zap.Logger
	users crud.UsersCrud
}

// NewUsersSrv constructs a UsersSrv.
func NewUsersSrv(log *zap.Logger, users crud.UsersCrud) UsersSrv {
	return &usersSrv{log: log, users: users}
}

// GetByAuthUserID returns the profile for an authenticated Supabase user.
func (s *usersSrv) GetByAuthUserID(ctx context.Context, authUserID uuid.UUID) (*model.User, *model.Err) {
	user, err := s.users.GetByAuthUserID(ctx, authUserID)
	if err != nil {
		var domainErr *model.Err
		if errors.As(err, &domainErr) {
			return nil, domainErr
		}
		s.log.Error("get user by auth id failed", zap.Error(err))
		return nil, constants.ErrInternal
	}
	return user, nil
}
