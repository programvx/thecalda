package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/programvx/thecalda/backend/internal/constants"
	"github.com/programvx/thecalda/backend/internal/middlewares"
	"github.com/programvx/thecalda/backend/internal/model"
	"github.com/programvx/thecalda/backend/internal/services"
)

// stubUsersSrv is a hand-written test double for services.UsersSrv.
type stubUsersSrv struct {
	user *model.User
	err  *model.Err
}

func (s *stubUsersSrv) GetByAuthUserID(_ context.Context, _ uuid.UUID) (*model.User, *model.Err) {
	return s.user, s.err
}

// newMeRouter builds a router for GET /api/me. When authUserID is non-nil it is
// injected into the context, standing in for the real auth middleware.
func newMeRouter(srv services.UsersSrv, authUserID *uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewUsersHandler(services.NewApiSrv(), srv)
	r.GET("/api/me", func(ctx *gin.Context) {
		if authUserID != nil {
			middlewares.SetAuthUserID(ctx, *authUserID)
		}
		ctx.Next()
	}, h.Me)
	return r
}

func TestUsersHandler_Me_Unauthenticated(t *testing.T) {
	router := newMeRouter(&stubUsersSrv{}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestUsersHandler_Me_Success(t *testing.T) {
	authUserID := uuid.New()
	want := &model.User{
		Base: model.Base{
			UID:       uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		AuthUserID: authUserID,
		Email:      "alice@example.com",
		FullName:   "Alice",
	}
	router := newMeRouter(&stubUsersSrv{user: want}, &authUserID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Data model.User `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Data.Email != want.Email {
		t.Fatalf("expected email %q, got %q", want.Email, resp.Data.Email)
	}
	if resp.Data.UID != want.UID {
		t.Fatalf("expected uid %s, got %s", want.UID, resp.Data.UID)
	}
}

func TestUsersHandler_Me_NotFound(t *testing.T) {
	authUserID := uuid.New()
	router := newMeRouter(&stubUsersSrv{err: constants.ErrNotFound}, &authUserID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
