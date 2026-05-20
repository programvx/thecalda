// Command api is the HTTP entrypoint for the thecalda backend.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/programvx/thecalda/backend/internal/core"
	"github.com/programvx/thecalda/backend/internal/db"
	"github.com/programvx/thecalda/backend/internal/db/crud"
	"github.com/programvx/thecalda/backend/internal/handlers"
	"github.com/programvx/thecalda/backend/internal/middlewares"
	"github.com/programvx/thecalda/backend/internal/routers"
	"github.com/programvx/thecalda/backend/internal/services"
	"github.com/programvx/thecalda/backend/internal/settings"
)

func main() {
	cfg := settings.NewSettings()

	logger, err := core.NewLogger(cfg.IsDev())
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	ctx := context.Background()

	store, err := db.NewStore(ctx, logger, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("database init failed", zap.Error(err))
	}
	defer func() { _ = store.Close() }()

	// CRUD
	usersCrud := crud.NewUsersCrud(store)

	// Services.
	apiSrv := services.NewApiSrv()
	usersSrv := services.NewUsersSrv(logger, usersCrud)

	// Auth middleware (fetches the Supabase JWKS up front).
	auth, err := middlewares.NewAuthMiddleware(logger, apiSrv, cfg.SupabaseURL)
	if err != nil {
		logger.Fatal("auth middleware init failed", zap.Error(err))
	}

	// HTTP engine and global middleware.
	if !cfg.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middlewares.Logger(logger))
	engine.Use(middlewares.Security())
	engine.Use(middlewares.CORS(cfg))

	// Routes — handler registration lives in the routers package.
	healthHandler := handlers.NewHealthHandler(store)
	usersHandler := handlers.NewUsersHandler(apiSrv, usersSrv)

	api := engine.Group("/api")
	routers.NewHealthRouter(api, healthHandler)

	api.Use(auth.Verify())

	routers.NewUsersRouter(api, usersHandler)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("api listening", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}
}
