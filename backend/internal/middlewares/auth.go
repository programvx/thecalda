package middlewares

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/services"
)

// ctxAuthUserID is the gin context key holding the authenticated user's id.
const ctxAuthUserID = "auth.userID"

// AuthMiddleware verifies Supabase-issued JWTs against the project's JWKS.
type AuthMiddleware struct {
	log    *zap.Logger
	apiSrv *services.ApiSrv
	jwks   keyfunc.Keyfunc
}

// NewAuthMiddleware builds an AuthMiddleware, fetching the Supabase JWKS up front.
func NewAuthMiddleware(log *zap.Logger, apiSrv *services.ApiSrv, supabaseURL string) (*AuthMiddleware, error) {
	jwksURL := strings.TrimRight(supabaseURL, "/") + "/auth/v1/.well-known/jwks.json"

	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("load supabase jwks from %s: %w", jwksURL, err)
	}

	return &AuthMiddleware{log: log, apiSrv: apiSrv, jwks: jwks}, nil
}

// Verify rejects any request that lacks a valid Supabase access token.
func (m *AuthMiddleware) Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raw, err := bearerToken(ctx)
		if err != nil {
			m.apiSrv.SendError(ctx, constants.ErrUnauthorized)
			return
		}

		token, err := jwt.Parse(
			raw,
			m.jwks.Keyfunc,
			jwt.WithValidMethods([]string{"ES256", "RS256"}),
			jwt.WithAudience("authenticated"),
			jwt.WithExpirationRequired(),
		)
		if err != nil || !token.Valid {
			m.log.Debug("jwt verification failed", zap.Error(err))
			m.apiSrv.SendError(ctx, constants.ErrUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.apiSrv.SendError(ctx, constants.ErrUnauthorized)
			return
		}

		sub, _ := claims["sub"].(string)
		authUserID, err := uuid.Parse(sub)
		if err != nil {
			m.apiSrv.SendError(ctx, constants.ErrUnauthorized)
			return
		}

		SetAuthUserID(ctx, authUserID)
		ctx.Next()
	}
}

// SetAuthUserID stores the authenticated user's Supabase id on the request context.
func SetAuthUserID(ctx *gin.Context, id uuid.UUID) {
	ctx.Set(ctxAuthUserID, id)
}

// AuthUserID returns the authenticated user's Supabase id from the request context.
func AuthUserID(ctx *gin.Context) (uuid.UUID, bool) {
	v, ok := ctx.Get(ctxAuthUserID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

func bearerToken(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	scheme, token, found := strings.Cut(header, " ")
	token = strings.TrimSpace(token)
	if !found || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", errors.New("malformed authorization header")
	}

	return token, nil
}
